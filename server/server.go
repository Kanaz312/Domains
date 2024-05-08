package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"goobl/game"
	"goobl/scenario"
	"log"
	"math/rand/v2"
	"net/http"
	"slices"
	"strconv"
	"text/template"
	"time"
)

type user struct {
	State game.Game
	Token int64
}

const maxUsers = 1000000

type ServerState struct {
	Users []user;
	Scenarios []scenario.Scenario;
	NumUsers int;
}

type sessionCreationError struct{
}

func (e *sessionCreationError) Error() string {
	return "Failed to create new session"
}

var sessionTokenName = "session"

func CalculateTomorrow() time.Time {
		now := time.Now()
		yyyy, mm, dd := now.Date()
		return time.Date(yyyy, mm, dd+1, 0, 0, 0, 0, now.Location())
}

func (s *ServerState) makeSession() (http.Cookie, error) {
	if s.NumUsers >= maxUsers {
		return http.Cookie{}, &sessionCreationError{}
	} else {
		token := s.Users[s.NumUsers].Token
		cookieValue := fmt.Sprintf("%d", token)
		tomorrow := CalculateTomorrow()
		secondsToTomorrow := int(time.Until(tomorrow).Seconds())

		newSessionCookie := http.Cookie{
			Name:       sessionTokenName,
			Value:      cookieValue,
			Expires:    tomorrow,
			MaxAge:     secondsToTomorrow,
			Secure:     true,
			HttpOnly:   true,
			SameSite:   http.SameSiteLaxMode,
		}

		s.NumUsers++

		return newSessionCookie, nil
	}
}

func (s *ServerState) deleteSession(w http.ResponseWriter) {
	deleteCookie := http.Cookie {
		Name: sessionTokenName,
		Value: "",
		Expires:    time.Now(),
		MaxAge: -1,
		Secure: true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &deleteCookie)
}

var indexTpl = template.Must(template.ParseFiles("index.html"))

func (s *ServerState) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionTokenName); err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			if newSessionCookie, err := s.makeSession(); err != nil {
				log.Println(err)
				return
			} else {
				http.SetCookie(w, &newSessionCookie)
			}
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	} else {
		if _, err := strconv.ParseInt(cookie.Value, 10, 64); err != nil {
			s.deleteSession(w)
			w.WriteHeader(http.StatusAccepted)
			w.Header().Add("HX-Redirect", "/")
			return
		}
	}
	
	if err := indexTpl.Execute(w, nil); err != nil {
		fmt.Printf("Failed to execute index.html %v\n", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

type sessionParsingError struct{
}

func (e *sessionParsingError) Error() string {
	return "Failed to parse session cookie"
}

func (s *ServerState) findUser(token int64) int {
	return slices.IndexFunc(s.Users[:], func(u user) bool {
		return u.Token == token
	})
}

func (s *ServerState) getUserSession(w http.ResponseWriter, r *http.Request) (int, error) {
	if cookie, err := r.Cookie(sessionTokenName); err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return -1, err
	} else {
		if token, err := strconv.ParseInt(cookie.Value, 10, 64); err != nil {
			s.deleteSession(w)
			return -1, &sessionParsingError{}
		} else {
			index := s.findUser(token)
			if index == -1 {
				s.deleteSession(w)
				return -1, &sessionParsingError{}
			} else {
				return index, nil
			}
		}
	}
}

type gameStateElements struct{
	Magic int;
	Population int;
	Sword int;
	Money int;
	Prompt string;
	Image string;
	LeftDescription string;
	RightDescription string;
	MagicLeftIndicator string;
	MagicRightIndicator string;
	PopulationLeftIndicator string;
	PopulationRightIndicator string;
	SwordLeftIndicator string;
	SwordRightIndicator string;
	MoneyLeftIndicator string;
	MoneyRightIndicator string;
	ReloadEndpoint string;
}

var gameStateElementsTpl = template.Must(template.ParseFiles("assets/static/gameStateElements.html"))

func (s *ServerState) populateDeadElements(state *game.Game, data *gameStateElements) error {
	scene := state.GetDeathScenario()
	if scene == nil {
		return errors.New("Failed to get death scenario")
	}

	data.Magic = state.Magic
	data.Population = state.Population
	data.Sword = state.Sword
	data.Money = state.Money
	data.Prompt = scene.Prompt
	data.Image = scene.Image
	data.LeftDescription = scene.Decisions[0].Description
	data.RightDescription = scene.Decisions[1].Description
	data.MagicLeftIndicator = getIndicatorImage(0)
	data.MagicRightIndicator = getIndicatorImage(0)
	data.PopulationLeftIndicator = getIndicatorImage(0)
	data.PopulationRightIndicator = getIndicatorImage(0)
	data.SwordLeftIndicator = getIndicatorImage(0)
	data.SwordRightIndicator = getIndicatorImage(0)
	data.MoneyLeftIndicator = getIndicatorImage(0)
	data.MoneyRightIndicator = getIndicatorImage(0)
	data.ReloadEndpoint = "results"

	return nil
}

func getIndicatorImage(ChangeInStat int) string {
	switch {
	case ChangeInStat > 20:
		return "BigPositive"
	case ChangeInStat > 0:
		return "SmallPositive"
	case ChangeInStat < -20:
		return "BigNegative"
	case ChangeInStat < 0:
		return "SmallNegative"
	default:
		return "Neutral"
	}
}

func (s *ServerState) populateLivingElements(state *game.Game, data *gameStateElements) {
	scene := &s.Scenarios[state.ScenarioIndex]

	data.Magic = state.Magic
	data.Population = state.Population
	data.Sword = state.Sword
	data.Money = state.Money
	data.Prompt = scene.Prompt
	data.Image = scene.Image
	data.LeftDescription = scene.Decisions[0].Description
	data.RightDescription = scene.Decisions[1].Description
	data.MagicLeftIndicator = getIndicatorImage(scene.Decisions[0].Magic)
	data.MagicRightIndicator = getIndicatorImage(scene.Decisions[1].Magic)
	data.PopulationLeftIndicator = getIndicatorImage(scene.Decisions[0].Population)
	data.PopulationRightIndicator = getIndicatorImage(scene.Decisions[1].Population)
	data.SwordLeftIndicator = getIndicatorImage(scene.Decisions[0].Sword)
	data.SwordRightIndicator = getIndicatorImage(scene.Decisions[1].Sword)
	data.MoneyLeftIndicator = getIndicatorImage(scene.Decisions[0].Money)
	data.MoneyRightIndicator = getIndicatorImage(scene.Decisions[1].Money)
	data.ReloadEndpoint = "gameStateElements"
}

func (s *ServerState) GameStateElementsHandler(w http.ResponseWriter, r *http.Request) {
	if index, err := s.getUserSession(w, r); err == nil {
		state := &s.Users[index].State
		data := gameStateElements{}

		if state.IsDead() {
			if err = s.populateDeadElements(state, &data); err != nil {
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
		} else {
			s.populateLivingElements(state, &data)
		}

		err := gameStateElementsTpl.Execute(w, data)
		if err != nil {
			fmt.Printf("Failed to execute assets/static/gameStateElements.html %v\n", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	} else {
		w.Header().Add("HX-Redirect", "/")
	}
}

type decisionRequest struct {
	Choice int `json:"decision"`
}

func (s *ServerState) DecisionHandler(w http.ResponseWriter, r *http.Request) {
	index, err := s.getUserSession(w, r)
	if err != nil {
		return
	}
	currentGameState := &s.Users[index].State
	currentScenario := s.Scenarios[currentGameState.ScenarioIndex]

	if currentGameState.IsDead() {
		return
	}
	
	var choice decisionRequest
	json.NewDecoder(r.Body).Decode(&choice)

	decisionIndex := 0
	if (choice.Choice == -1) {
		decisionIndex = 0
		w.Write([]byte("200 - Left Received"))
	} else if (choice.Choice == 1) {
		decisionIndex = 1
		w.Write([]byte("200 - Right Received"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Uknown decision value"))
		return
	}

	currentGameState.ApplyDecision(&currentScenario.Decisions[decisionIndex], s.Scenarios)
}

var resultsTpl = template.Must(template.ParseFiles("assets/static/results.html"))


type resultsElements struct{
	CardsSeen int;
}

func (s *ServerState) ResultsHandler(w http.ResponseWriter, r *http.Request) {
	index, err := s.getUserSession(w, r)
	if err != nil {
		return
	}
	currentGameState := &s.Users[index].State

	data := resultsElements{currentGameState.ScenarioIndex}

	err = resultsTpl.Execute(w, data)
	if err != nil {
		fmt.Printf("Failed to execute assets/static/results.html %v\n", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func (s *ServerState) InitializeServer() {
	s.Users = make([]user, maxUsers)
	for i := 0; i < len(s.Users); i++ {
		u := &s.Users[i]
		u.State = game.Game{Magic: 50, Population: 50, Sword: 50, Money: 50, ScenarioIndex: 0}
		u.Token = rand.Int64()
	}
	
	s.Scenarios = scenario.Scenarios
	rand.Shuffle(len(s.Scenarios), func(i, j int) {
		s.Scenarios[i], s.Scenarios[j] = s.Scenarios[j], s.Scenarios[i]
	})

	startScenario := scenario.Scenario{
		Prompt: "I heard you've been feeling down. I'll throw a show for you!",
		Image: "BaseGoobl",
		Decisions: 
		[2]scenario.Decision{
			{Magic: 0, Population: 0, Sword: 0, Money: 0, Description: "Awww, thanks!"},
			{Magic: 0, Population: 0, Sword: 0, Money: 0, Description: "I'm guessing I have no choice..."}},
		}

	s.Scenarios = slices.Insert(s.Scenarios, 0, startScenario)

	s.NumUsers = 0
}
