package main

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	serverMiddleware "github.com/djokcik/praktikum-go-devops/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	cfg := server.NewConfig()

	logging.
		NewLogger().
		Info().
		Msgf("config: %+v", cfg)

	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)
	mux.Use(serverMiddleware.GzipHandle)
	mux.Use(serverMiddleware.LoggerMiddleware())

	makeMetricRoutes(ctx, wg, mux, cfg)

	go func() {
		err := http.ListenAndServe(cfg.Address, mux)
		if err != nil {
			logging.NewLogger().Fatal().Err(err).Msg("server stopped")
		}

	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	cancel()
	logging.NewLogger().Info().Msg("Shutdown Server ...")
	wg.Wait()
}
