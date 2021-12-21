package main

import (
	"fmt"
	"github.com/Jokcik/praktikum-go-devops/internal/server/handler"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage/model"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func makeMetricRoutes(mux *chi.Mux) *handler.Handler {
	rr := storage.NewRepositoryRegistry(new(model.Database), &storage.MetricRepository{})
	metricRepository, err := rr.Repository("MetricRepository")
	if err != nil {
		fmt.Println("Error provide repository 'MetricRepository'")
	}

	h := handler.NewHandler(mux, metricRepository)

	h.Get("/", h.ListHandler())

	h.Route("/update", func(r chi.Router) {
		r.Post("/counter/{name}/{value}", h.CounterHandler())
		r.Post("/gauge/{name}/{value}", h.GaugeHandler())
		r.Post("/counter/*", http.NotFound)
		r.Post("/gauge/*", http.NotFound)
		r.Post("/*", handler.NotImplementedHandler)
	})

	h.Route("/value", func(r chi.Router) {
		r.Get("/counter/{name}", h.GetCounterHandler())
		r.Get("/gauge/{name}", h.GetGaugeHandler())
		r.Get("/*", http.NotFound)
	})

	return h
}
