// Package handler is a collection of handlers for use get and save metric.
// The packages includes handlers which update and get metrics in json and row format.
package handler

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/service"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/reporegistry"
	commonService "github.com/djokcik/praktikum-go-devops/internal/service"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// Handler struct for all handlers and require DI dependencies
type Handler struct {
	*chi.Mux
	Hash    commonService.HashService
	Counter service.CounterService
	Gauge   service.GaugeService
}

// NewHandler constructor for Handler
func NewHandler(mux *chi.Mux, cfg server.Config, repoRegistry reporegistry.RepoRegistry) *Handler {
	hashService := commonService.NewHashService(cfg.Key)

	return &Handler{
		Mux:     mux,
		Hash:    hashService,
		Counter: &service.CounterServiceImpl{Repo: repoRegistry.GetCounterRepo(), Hash: hashService},
		Gauge:   &service.GaugeServiceImpl{Repo: repoRegistry.GetGaugeRepo(), Hash: hashService},
	}
}

func (h *Handler) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "Handler").Logger()

	return &logger
}
