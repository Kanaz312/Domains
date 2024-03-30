package main

import (
	"fmt"
	"goobl/game"
	"goobl/scenarios"
	"goobl/server"
	"net/http"
)


func main() {
	left0 := scenario.Decision{Cross: 30, Population: 10, Sword: 0, Money: -40, Description: "Build the church"}
	right0 := scenario.Decision{Cross: -20, Population: 0, Sword: 0, Money: 10, Description: "No good, costs too much"}
	scenario0 := scenario.Scenario{Prompt: "We would like to construct a cathedral to spread the Good Word", Image: "jack", LeftDecision: left0, RightDecision: right0}
	left1 := scenario.Decision{Cross: -10, Population: 40, Sword: 0, Money: -10, Description: "Clear the forests to make new farmland"}
	right1 := scenario.Decision{Cross: 0, Population: 30, Sword: -30, Money: 0, Description: "Steal from our neighbors"}
	scenario1 := scenario.Scenario{Prompt: "The people are going hungry", Image: "queen", LeftDecision: left1, RightDecision: right1}

	mainServerState := server.ServerState{GameStates: make([]game.Game, 0, 5), Scenarios: make([]scenario.Scenario, 2, 5), CookieIndex: 0}

	mainServerState.Scenarios[0] = scenario0
	mainServerState.Scenarios[1] = scenario1

	fs := http.FileServer(http.Dir("assets"))

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", mainServerState.IndexHandler)
	http.HandleFunc("/gameStateElements", mainServerState.GameStateElementsHandler)
	http.HandleFunc("/decide", mainServerState.DecisionHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to listen and start %v\n", err)
	}
}
