package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	var cfg server.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	makeMetricRoutes(mux, &cfg)

	log.Fatal(http.ListenAndServe(cfg.Address, mux))
}
