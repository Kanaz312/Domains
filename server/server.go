package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"goobl/game"
	"goobl/scenarios"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

type ServerState struct {
	GameStates []game.Game;
	Scenarios []scenario.Scenario;
	CookieIndex int;
}

func (s *ServerState) makeSession() (int, http.Cookie) {
	index := s.CookieIndex
	cookieValue := fmt.Sprintf("%d", index)

	newSessionCookie := http.Cookie {
		Name: "session",
		Value: cookieValue,

		Secure: true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	sampleState := game.Game{Cross: 50, Population: 50, Sword: 50, Money: 50, Scenario: s.Scenarios[0], ScenarioIndex: 0}
	s.GameStates = append(s.GameStates, sampleState)

	s.CookieIndex++

	return index, newSessionCookie
}

var indexTpl = template.Must(template.ParseFiles("index.html"))

func (s *ServerState) IndexHandler(w http.ResponseWriter, r *http.Request) {
	index := 0
	if cookie, err := r.Cookie("session"); err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			newIndex, newSessionCookie := s.makeSession()
			index = newIndex
			http.SetCookie(w, &newSessionCookie)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	} else {
		index, err = strconv.Atoi(cookie.Value)
		if err != nil || index >= s.CookieIndex {
			http.Error(w, "invalid session", http.StatusUnauthorized)
			return
		}
	}
	
	if err := indexTpl.Execute(w, s.GameStates[index]); err != nil {
		fmt.Printf("Failed to execute index.html %v\n", err)
	}
}

type sessionParsingError struct{
}

func (e *sessionParsingError) Error() string {
	return "Failed to parse session cookie"
}

func (s *ServerState) getUserSession(w http.ResponseWriter, r *http.Request) (int, error) {
	if cookie, err := r.Cookie("session"); err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "session token not found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}

		return -1, err
	} else {
		index, err := strconv.Atoi(cookie.Value)
		if err != nil || index >= s.CookieIndex {
			http.Error(w, "invalid session", http.StatusUnauthorized)
			return -1, &sessionParsingError{}
		} else {
			return index, nil
		}
	}
}

var gameStateElementsTpl = template.Must(template.ParseFiles("assets/static/gameStateElements.html"))

func (s *ServerState) GameStateElementsHandler(w http.ResponseWriter, r *http.Request) {
	if index, err := s.getUserSession(w, r); err == nil {
		err := gameStateElementsTpl.Execute(w, s.GameStates[index])
		if err != nil {
			fmt.Printf("Failed to execute assets/static/gameStateElements.html %v\n", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
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
	currentGameState := &s.GameStates[index]
	currentScenario := currentGameState.Scenario
	
	var choice decisionRequest
	json.NewDecoder(r.Body).Decode(&choice)

	if (choice.Choice == -1) {
		currentGameState.Cross += currentScenario.LeftDecision.Cross
		currentGameState.Population += currentScenario.LeftDecision.Population
		currentGameState.Sword += currentScenario.LeftDecision.Sword
		currentGameState.Money += currentScenario.LeftDecision.Money

		currentGameState.ScenarioIndex = (currentGameState.ScenarioIndex + 1) % len(s.Scenarios)
		currentGameState.Scenario = s.Scenarios[currentGameState.ScenarioIndex]

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("200 - Left Received"))
	} else if (choice.Choice == 1) {
		currentGameState.Cross += currentScenario.RightDecision.Cross
		currentGameState.Population += currentScenario.RightDecision.Population
		currentGameState.Sword += currentScenario.RightDecision.Sword
		currentGameState.Money += currentScenario.RightDecision.Money

		currentGameState.ScenarioIndex = (currentGameState.ScenarioIndex + 1) % len(s.Scenarios)
		currentGameState.Scenario = s.Scenarios[currentGameState.ScenarioIndex]

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("200 - Right Received"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Uknown decision value"))
	}
}
