package main

import (
	"context"
	"goobl/server"
	"log"
	"net/http"
	"time"
)

func startServer(goServer *http.Server){
	if err := goServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe failed: %v", err)
	}

}

func shutdownServer(goServer *http.Server){
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := goServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
}

func main() {
	mainServerState := server.ServerState{}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", mainServerState.IndexHandler)
	mux.HandleFunc("/gameStateElements", mainServerState.GameStateElementsHandler)
	mux.HandleFunc("/decide", mainServerState.DecisionHandler)
	mux.HandleFunc("/results", mainServerState.ResultsHandler)

	for {
		mainServerState.InitializeServer()

		goServer := &http.Server{
			Addr:                         ":8080",
			Handler:                      mux,
			DisableGeneralOptionsHandler: false,
			ReadTimeout:                  10 * time.Second,
			ReadHeaderTimeout:            10 * time.Second,
			WriteTimeout:                 10 * time.Second,
			IdleTimeout:                  0,
			MaxHeaderBytes:               1 << 20,
		}

		go startServer(goServer)

		tomorrow := server.CalculateTomorrow()
		time.Sleep(time.Until(tomorrow))

		shutdownServer(goServer)

		log.Print("Re-initializing server...")
	}
}
