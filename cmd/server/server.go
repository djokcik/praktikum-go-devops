package main

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/handler"
	"github.com/djokcik/praktikum-go-devops/internal/server/service"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/reporegistry"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"sync"
)

func makeMetricRoutes(ctx context.Context, wg *sync.WaitGroup, mux *chi.Mux, cfg server.Config) *handler.Handler {
	var repoRegistry reporegistry.RepoRegistry

	if cfg.DatabaseDsn != "" {
		databaseService := service.NewDatabaseService(ctx, cfg)
		db, err := databaseService.Open(ctx)
		if err != nil {
			logging.NewLogger().Fatal().Err(err).Msgf("Doesn`t open database connection")
			os.Exit(1)
		}

		repoRegistry = reporegistry.NewPostgreSQL(ctx, db)
	} else {
		repoRegistry = reporegistry.NewInMem(ctx, wg, cfg)
	}

	h := handler.NewHandler(mux, cfg, repoRegistry)

	h.Get("/", h.ListHandler())
	h.Get("/ping", h.PingHandler(repoRegistry))

	h.Route("/updates", func(r chi.Router) {
		r.Post("/", h.UpdateListJSONHandler())
	})

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
