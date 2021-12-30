package handler

import (
	"github.com/djokcik/praktikum-go-devops/internal/server/service"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/go-chi/chi/v5"
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
