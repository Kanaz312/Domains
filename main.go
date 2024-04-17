package main

import (
	"fmt"
	"goobl/server"
	"net/http"
)


func main() {
	mainServerState := server.ServerState{}
	mainServerState.InitializeServer()

	fs := http.FileServer(http.Dir("assets"))

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", mainServerState.IndexHandler)
	http.HandleFunc("/gameStateElements", mainServerState.GameStateElementsHandler)
	http.HandleFunc("/decide", mainServerState.DecisionHandler)
	http.HandleFunc("/results", mainServerState.ResultsHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to listen and start %v\n", err)
	}
}
