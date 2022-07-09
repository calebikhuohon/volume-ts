package main

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"volume-ts/pkg/handler"
	"volume-ts/pkg/service"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		systemCall := <-exit
		log.Printf("System call: %+v", systemCall)
		cancel()
	}()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("Failed to parse app config: %s", err)
	}
	if err := config.Validate(); err != nil {
		log.Fatalf("Failed to validate config: %s", err)
	}

	flightService := service.NewService()

	r := handler.New(flightService, handler.Config{Timeout: 10 * time.Second})

	srv := NewServer(config.Ports.HTTP, r)

	go func() {
		log.Print("Starting http server..")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Could not listen and serve: %s", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	log.Print("Stopping the http server..")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Could not shut down gracefully: %s", err)
		os.Exit(1)
	}

	defer cancel()
}

func NewServer(addr string, h http.Handler) *http.Server {
	s := &http.Server{Addr: addr, Handler: h}

	return s
}
