package handler_test

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/handler"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/reporegistry"
	"github.com/go-chi/chi/v5"
	"sync"
)

func Example() {
	mux := chi.NewMux()
	cfg := server.Config{}

	repoRegistry := reporegistry.NewInMem(context.Background(), new(sync.WaitGroup), cfg)

	h := handler.NewHandler(mux, cfg, repoRegistry)

	h.Get("/", h.ListHandler())
	h.Get("/ping", h.PingHandler(repoRegistry))
}
