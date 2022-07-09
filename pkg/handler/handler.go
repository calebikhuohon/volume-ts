package handler

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

type FlightService interface {
	SortFlightPaths(ctx context.Context, user string) ([]string, error)
}

type Handler struct {
	router        http.Handler
	flightService FlightService
}

type Config struct {
	Timeout time.Duration
}

func New(
	flightService FlightService,
	config Config,
) http.Handler {
	r := chi.NewRouter()

	h := &Handler{
		router:        r,
		flightService: flightService,
	}

	timeout := 10 * time.Second
	if config.Timeout > 0 {
		timeout = config.Timeout
	}

	r.Use(
		chimiddleware.Timeout(timeout),
		chimiddleware.SetHeader("Content-Type", "application/json"),
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			AllowedHeaders:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "HEAD", "OPTIONS", "PUT"},
			Debug:            false,
		}).Handler,
	)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	r.Route("/sort-flights", func(r chi.Router) {
		r.Get("/", h.SortFlights)
	})

	return h
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
