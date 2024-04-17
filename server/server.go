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

func (s *ServerState) makeSession() (http.Cookie, error) {
	if s.NumUsers >= maxUsers {
		return http.Cookie{}, &sessionCreationError{}
	} else {
		token := s.Users[s.NumUsers].Token
		cookieValue := fmt.Sprintf("%d", token)

		newSessionCookie := http.Cookie {
			Name: sessionTokenName,
			Value: cookieValue,

			Secure: true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}

		s.NumUsers++

		return newSessionCookie, nil
	}
}

func (s *ServerState) deleteSession(w http.ResponseWriter) {
	deleteCookie := http.Cookie {
		Name: sessionTokenName,
		Value: "",

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
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "session token not found", http.StatusBadRequest)
		default:
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
}

var deadTpl = template.Must(template.ParseFiles("assets/static/dead.html"))
var gameStateElementsTpl = template.Must(template.ParseFiles("assets/static/gameStateElements.html"))

func (s *ServerState) populateDeadElements(w http.ResponseWriter, state *game.Game) {
	scene := state.GetDeathScenario()
	if scene == nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	data := gameStateElements{
		state.Magic,
		state.Population,
		state.Sword,
		state.Money,
		scene.Prompt,
		scene.Image,
		scene.Decisions[0].Description,
		scene.Decisions[1].Description,
	}

	if err:= deadTpl.Execute(w, data); err != nil {
		fmt.Printf("Failed to execute assets/static/dead.html %v\n", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func (s *ServerState) populateLivingElements(w http.ResponseWriter, state *game.Game) {
	scene := &s.Scenarios[state.ScenarioIndex]

	data := gameStateElements{
		state.Magic,
		state.Population,
		state.Sword,
		state.Money,
		scene.Prompt,
		scene.Image,
		scene.Decisions[0].Description,
		scene.Decisions[1].Description,
	}

	err := gameStateElementsTpl.Execute(w, data)
	if err != nil {
		fmt.Printf("Failed to execute assets/static/gameStateElements.html %v\n", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

}

func (s *ServerState) GameStateElementsHandler(w http.ResponseWriter, r *http.Request) {
	if index, err := s.getUserSession(w, r); err == nil {
		state := &s.Users[index].State
		if state.IsDead() {
			s.populateDeadElements(w, state)
		} else {
			s.populateLivingElements(w, state)
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

func (s *ServerState) ResultsHandler(w http.ResponseWriter, r *http.Request) {
	err := resultsTpl.Execute(w, nil)
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
