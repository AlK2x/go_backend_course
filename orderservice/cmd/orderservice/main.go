package main

import (
	"context"
	"log"
	"net/http"
	"orderservice/pkg/orderservice"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	router := orderservice.NewRouter()
	server := &http.Server{
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8000",
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	waitForShutdown(server)
}

func waitForShutdown(server *http.Server) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT)
	<-sigint

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Shutdown server error: %v", err)
	}
}
