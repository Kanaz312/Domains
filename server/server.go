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

var maxUsers = 10000

type ServerState struct {
	Users [10000]user;
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
	Cross int;
	Population int;
	Sword int;
	Money int;
	Prompt string;
	Image string;
	LeftDescription string;
	RightDescription string;
}

var gameStateElementsTpl = template.Must(template.ParseFiles("assets/static/gameStateElements.html"))

func (s *ServerState) GameStateElementsHandler(w http.ResponseWriter, r *http.Request) {
	if index, err := s.getUserSession(w, r); err == nil {
		state := s.Users[index].State
		scenario := s.Scenarios[state.ScenarioIndex]
		data := gameStateElements{
			state.Cross,
			state.Population,
			state.Sword,
			state.Money,
			scenario.Prompt,
			scenario.Image,
			scenario.Decisions[0].Description,
			scenario.Decisions[1].Description}
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
	
	var choice decisionRequest
	json.NewDecoder(r.Body).Decode(&choice)

	decisionIndex := 0
	if (choice.Choice == -1) {
		decisionIndex = 0

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("200 - Left Received"))
	} else if (choice.Choice == 1) {
		decisionIndex = 1

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("200 - Right Received"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Uknown decision value"))
		return
	}

	decision := currentScenario.Decisions[decisionIndex];
	currentGameState.Cross += decision.Cross
	currentGameState.Population += decision.Population
	currentGameState.Sword += decision.Sword
	currentGameState.Money += decision.Money

	currentGameState.ScenarioIndex = (currentGameState.ScenarioIndex + 1) % len(s.Scenarios)
}

func (s *ServerState) InitializeServer() {
	for i := 0; i < len(s.Users); i++ {
		u := &s.Users[i]
		u.State = game.Game{Cross: 50, Population: 50, Sword: 50, Money: 50, ScenarioIndex: 0}
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
			{Cross: 0, Population: 0, Sword: 0, Money: 0, Description: "Awww, thanks!"},
			{Cross: 0, Population: 0, Sword: 0, Money: 0, Description: "I'm guessing I have no choice..."}},
		}

	s.Scenarios = slices.Insert(s.Scenarios, 0, startScenario)

	s.NumUsers = 0
}
