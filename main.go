package main

import (
	"fmt"
	"goobl/game"
	"goobl/scenario"
	"goobl/server"
	"net/http"
)


func main() {
	mainServerState := server.ServerState{GameStates: make([]game.Game, 0, 5), Scenarios: scenario.Scenarios, CookieIndex: 0}

	fs := http.FileServer(http.Dir("assets"))

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", mainServerState.IndexHandler)
	http.HandleFunc("/gameStateElements", mainServerState.GameStateElementsHandler)
	http.HandleFunc("/decide", mainServerState.DecisionHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to listen and start %v\n", err)
	}
}
