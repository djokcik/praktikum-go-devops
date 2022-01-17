package handler

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/server/service"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Handler struct {
	*chi.Mux
	Counter service.CounterService
	Gauge   service.GaugeService
}

func NewHandler(mux *chi.Mux, repo storage.Repository) *Handler {
	return &Handler{
		Mux:     mux,
		Counter: &service.CounterServiceImpl{Repo: repo},
		Gauge:   &service.GaugeServiceImpl{Repo: repo},
	}
}

func (h *Handler) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "Handler").Logger()

	return &logger
}
