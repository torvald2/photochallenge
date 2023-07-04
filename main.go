package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/torvald2/photochallenge/handlers"
	"github.com/torvald2/photochallenge/recognizer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	fr, err := recognizer.NewTokenizer(0.6)
	if err != nil {
		panic(err)
	}
	router := handlers.NewRouter(fr)

	srv := &http.Server{
		Addr:    ":8085",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()
	srv.Shutdown(ctx)
	os.Exit(0)
}
