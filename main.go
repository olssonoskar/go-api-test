package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/olsson/my-api/handlers"
)

// Simple API based on Nic Jacksson's 'Building Microservices with Go'

func main() {
	fmt.Println("Server starting...")

	l := log.New(os.Stdout, "my-api", log.LstdFlags)
	albumh := handlers.NewAlbums(l)
	simpleh := handlers.NewSimple(l)

	mux := http.NewServeMux()
	mux.Handle("/albums/", albumh)
	mux.Handle("/hello", simpleh)

	serv := &http.Server{
		Addr:         ":9090",
		Handler:      mux,
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
