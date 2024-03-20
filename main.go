package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type gameState struct {
	Cross int;
	Population int;
	Sword int;
	Money int;
}

type serverState struct {
	gameStates []gameState;
}


var indexTpl = template.Must(template.ParseFiles("index.html"))

func (s *serverState) indexHandler(w http.ResponseWriter, r *http.Request) {
	err := indexTpl.Execute(w, s.gameStates[0])
	if err != nil {
		fmt.Printf("Failed to execute index.html %v\n", err)
	}
}

var statTpl = template.Must(template.ParseFiles("assets/static/stats.html"))

func (s *serverState) statsHandler(w http.ResponseWriter, r *http.Request) {
	err := statTpl.Execute(w, s.gameStates[0])
	if err != nil {
		fmt.Printf("Failed to execute assets/static/stats.html %v\n", err)
	}
}

func (s *serverState) leftHandler(w http.ResponseWriter, r *http.Request) {
	gameState := &s.gameStates[0]
	gameState.Cross += 20
	gameState.Population -= 20
	gameState.Sword += 20
	gameState.Money -= 20

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("200 - Left Received"))
	log.Printf("Handled left, new stats: %d %d %d %d\n", gameState.Cross, gameState.Population, gameState.Sword, gameState.Money)
}

func (s *serverState) rightHandler(w http.ResponseWriter, r *http.Request) {
	gameState := &s.gameStates[0]
	gameState.Cross -= 20
	gameState.Population += 20
	gameState.Sword -= 20
	gameState.Money += 20

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("200 - Right Received"))
	log.Printf("Handled right, new stats: %d %d %d %d\n", gameState.Cross, gameState.Population, gameState.Sword, gameState.Money)
}

func main() {
	mainServerState := serverState{make([]gameState, 5)}
	state := gameState{50, 50, 50, 50}
	mainServerState.gameStates[0] = state;

	fs := http.FileServer(http.Dir("assets"))

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", mainServerState.indexHandler)
	http.HandleFunc("/stats", mainServerState.statsHandler)
	http.HandleFunc("/left", mainServerState.leftHandler)
	http.HandleFunc("/right", mainServerState.rightHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to listen and start %v\n", err)
	}
}
