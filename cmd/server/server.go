package main

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/handler"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sync"
)

func makeMetricRoutes(ctx context.Context, wg *sync.WaitGroup, mux *chi.Mux, cfg server.Config) *handler.Handler {
	metricRepository := storage.NewMetricRepository(ctx, wg, cfg)

	h := handler.NewHandler(mux, cfg, metricRepository)

	h.Get("/", h.ListHandler())
	h.Get("/ping", h.PingHandler(metricRepository))

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
