package main

import (
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/handler"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/model"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func makeMetricRoutes(mux *chi.Mux, cfg *server.Config) *handler.Handler {
	rr := storage.NewRepositoryRegistry(cfg, new(model.Database), &storage.MetricRepository{})
	metricRepository, err := rr.Repository("MetricRepository")
	if err != nil {
		log.Println("Error provide repository 'MetricRepository'")
	}

	h := handler.NewHandler(mux, metricRepository)

	h.Get("/", h.ListHandler())

	h.Route("/update", func(r chi.Router) {
		r.Post("/", h.UpdateJSONHandler())
		r.Post("/counter/{name}/{value}", h.CounterHandler())
		r.Post("/gauge/{name}/{value}", h.GaugeHandler())
		r.Post("/counter/*", http.NotFound)
		r.Post("/gauge/*", http.NotFound)
		r.Post("/*", handler.NotImplementedHandler)
	})

	h.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetJSONHandler())
		r.Get("/counter/{name}", h.GetCounterMetricHandler())
		r.Get("/gauge/{name}", h.GetGaugeMetricHandler())
		r.Get("/*", http.NotFound)
	})

	return h
}
