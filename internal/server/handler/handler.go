package handler

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/service"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	commonService "github.com/djokcik/praktikum-go-devops/internal/service"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Handler struct {
	*chi.Mux
	Hash    commonService.HashService
	Counter service.CounterService
	Gauge   service.GaugeService
}

func NewHandler(mux *chi.Mux, cfg server.Config, repo storage.MetricRepository) *Handler {
	hashService := &commonService.HashServiceImpl{HashKey: cfg.Key}

	return &Handler{
		Mux:     mux,
		Hash:    hashService,
		Counter: &service.CounterServiceImpl{Repo: repo, Hash: hashService},
		Gauge:   &service.GaugeServiceImpl{Repo: repo, Hash: hashService},
	}
}

func (h *Handler) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "Handler").Logger()

	return &logger
}
