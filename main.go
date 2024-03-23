package main

import (
	"fmt"
	"html/template"
	"log"
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

var statTpl = template.Must(template.ParseFiles("assets/static/stats.html"))

func (s *serverState) statsHandler(w http.ResponseWriter, r *http.Request) {
	err := statTpl.Execute(w, s.GameStates[0])
	if err != nil {
		fmt.Printf("Failed to execute assets/static/stats.html %v\n", err)
	}
}

var scenarioTpl = template.Must(template.ParseFiles("assets/static/scenario.html"))

func (s *serverState) scenarioHandler(w http.ResponseWriter, r *http.Request) {
	err := scenarioTpl.Execute(w, s.GameStates[0].Scenario)
	if err != nil {
		fmt.Printf("Failed to execute assets/static/scenario.html %v\n", err)
	}
}

var gameStateElementsTpl = template.Must(template.ParseFiles("assets/static/gameStateElements.html"))

func (s *serverState) gameStateElementsHandler(w http.ResponseWriter, r *http.Request) {
	err := gameStateElementsTpl.Execute(w, s.GameStates[0])
	if err != nil {
		fmt.Printf("Failed to execute assets/static/gameStateElements.html %v\n", err)
	}
}

func (s *serverState) leftHandler(w http.ResponseWriter, r *http.Request) {
	currentGameState := &s.GameStates[0]
	currentScenario := currentGameState.Scenario
	
	currentGameState.Cross += currentScenario.LeftDecision.Cross
	currentGameState.Population += currentScenario.LeftDecision.Population
	currentGameState.Sword += currentScenario.LeftDecision.Sword
	currentGameState.Money += currentScenario.LeftDecision.Money

	currentGameState.ScenarioIndex = (currentGameState.ScenarioIndex + 1) % len(s.Scenarios)
	currentGameState.Scenario = s.Scenarios[currentGameState.ScenarioIndex]
	log.Printf("%v", currentGameState.Scenario)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("200 - Left Received"))
	log.Printf("Handled left, new stats: %d %d %d %d\n", currentGameState.Cross, currentGameState.Population, currentGameState.Sword, currentGameState.Money)
}

func (s *serverState) rightHandler(w http.ResponseWriter, r *http.Request) {
	currentGameState := &s.GameStates[0]
	currentScenario := currentGameState.Scenario
	
	currentGameState.Cross += currentScenario.RightDecision.Cross
	currentGameState.Population += currentScenario.RightDecision.Population
	currentGameState.Sword += currentScenario.RightDecision.Sword
	currentGameState.Money += currentScenario.RightDecision.Money

	currentGameState.ScenarioIndex = (currentGameState.ScenarioIndex + 1) % len(s.Scenarios)
	currentGameState.Scenario = s.Scenarios[currentGameState.ScenarioIndex]

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("200 - Right Received"))
	log.Printf("Handled right, new stats: %d %d %d %d\n", currentGameState.Cross, currentGameState.Population, currentGameState.Sword, currentGameState.Money)
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
	http.HandleFunc("/stats", mainServerState.statsHandler)
	http.HandleFunc("/scenario", mainServerState.scenarioHandler)
	http.HandleFunc("/gameStateElements", mainServerState.gameStateElementsHandler)
	http.HandleFunc("/left", mainServerState.leftHandler)
	http.HandleFunc("/right", mainServerState.rightHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to listen and start %v\n", err)
	}
}
