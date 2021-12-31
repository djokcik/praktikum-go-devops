package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type config struct {
	Address string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

func main() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	makeMetricRoutes(mux)

	log.Fatal(http.ListenAndServe(cfg.Address, mux))
}
