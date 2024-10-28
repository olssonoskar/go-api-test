package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/olsson/my-api/handlers"
)

// Simple API based on Nic Jacksson's 'Building Microservices with Go'

func main() {
	fmt.Println("Server starting...")

	l := log.New(os.Stdout, "my-api", log.LstdFlags)
	albumh := handlers.NewAlbums(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", albumh.GetAll)
	getRouter.HandleFunc("/{id:[0-9]+}", albumh.GetAlbum)

	putRouter := sm.Methods("PUT").Subrouter()
	putRouter.Use(albumh.ValidateAlbum)
	putRouter.HandleFunc("/{id:[0-9]+}", albumh.Update)

	postRouter := sm.Methods("POST").Subrouter()
	postRouter.Use(albumh.ValidateAlbum)
	postRouter.HandleFunc("/", albumh.Create)

	deleteRouter := sm.Methods("DELETE").Subrouter()
	deleteRouter.HandleFunc("/", albumh.Delete)

	serv := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		err := serv.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	schan := make(chan os.Signal, 2)
	signal.Notify(schan, os.Interrupt)

	sig := <-schan
	l.Println("Recieved signal to terminate:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	serv.Shutdown(ctx)
}
