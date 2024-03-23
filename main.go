package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type decision struct {
	Cross int;
	Population int;
	Sword int;
	Money int;
	Description string;

}

type scenario struct {
	Prompt string;
	Image string;
	LeftDecision decision;
	RightDecision decision;
}

type gameState struct {
	Cross int;
	Population int;
	Sword int;
	Money int;
	Scenario scenario;
	ScenarioIndex int;
}

type serverState struct {
	GameStates []gameState;
	Scenarios []scenario;
}


var indexTpl = template.Must(template.ParseFiles("index.html"))

func (s *serverState) indexHandler(w http.ResponseWriter, r *http.Request) {
	err := indexTpl.Execute(w, s.GameStates[0])
	if err != nil {
		fmt.Printf("Failed to execute index.html %v\n", err)
	}
}

var gameStateElementsTpl = template.Must(template.ParseFiles("assets/static/gameStateElements.html"))

func (s *serverState) gameStateElementsHandler(w http.ResponseWriter, r *http.Request) {
	err := gameStateElementsTpl.Execute(w, s.GameStates[0])
	if err != nil {
		fmt.Printf("Failed to execute assets/static/gameStateElements.html %v\n", err)
	}
}

type decisionRequest struct {
	Choice int `json:"decision"`
}

func (s *serverState) decisionHandler(w http.ResponseWriter, r *http.Request) {
	currentGameState := &s.GameStates[0]
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

func main() {
	left0 := decision{30, 10, 0, -40, "Build the church"}
	right0 := decision{-20, 0, 0, 10, "No good, costs too much"}
	scenario0 := scenario{"We would like to construct a cathedral to spread the Good Word", "jack", left0, right0}
	left1 := decision{-10, 40, 0, -10, "Clear the forests to make new farmland"}
	right1 := decision{0, 30, -30, 0, "Steal from our neighbors"}
	scenario1 := scenario{"The people are going hungry", "queen", left1, right1}

	mainServerState := serverState{make([]gameState, 5), make([]scenario, 2)}

	mainServerState.Scenarios[0] = scenario0
	mainServerState.Scenarios[1] = scenario1

	sampleState := gameState{50, 50, 50, 50, scenario0, 0}
	mainServerState.GameStates[0] = sampleState

	fs := http.FileServer(http.Dir("assets"))

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", mainServerState.indexHandler)
	http.HandleFunc("/gameStateElements", mainServerState.gameStateElementsHandler)
	http.HandleFunc("/decide", mainServerState.decisionHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to listen and start %v\n", err)
	}
}
