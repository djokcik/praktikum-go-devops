package handler

import (
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	*chi.Mux
	Repo storage.Repository
}

func NewHandler(mux *chi.Mux, repo storage.Repository) *Handler {
	return &Handler{
		Mux:  mux,
		Repo: repo,
	}
}
