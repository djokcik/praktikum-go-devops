package main

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	serverMiddleware "github.com/djokcik/praktikum-go-devops/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	cfg := server.NewConfig()
	log := logging.NewLogger()

	log.Info().Msgf("config: %+v", cfg)
	log.Info().Msgf("Build version: %s", buildVersion)
	log.Info().Msgf("Build date: %s", buildDate)
	log.Info().Msgf("Build commit: %s", buildCommit)

	repoRegistry := getRepoRegistry(ctx, cfg, wg)

	go func() {
		if cfg.GRPCAddress != "" {
			s, listen, err := makeGRPCMetricServer(cfg, repoRegistry)
			if err != nil {
				log.Fatal().Err(err).Msg("error listen gRPC server")
			}

			log.Info().Msgf("gRPC server started with address: %s", cfg.GRPCAddress)
			if err := s.Serve(listen); err != nil {
				log.Fatal().Err(err).Msg("error start gRPC server")
			}
		}
	}()

	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)
	mux.Use(serverMiddleware.TrustedSubnetHandle(cfg))
	mux.Use(serverMiddleware.GzipHandle)
	mux.Use(serverMiddleware.LoggerMiddleware())

	makeMetricRoutes(ctx, repoRegistry, mux, cfg)

	go func() {
		err := http.ListenAndServe(cfg.Address, mux)
		if err != nil {
			log.Fatal().Err(err).Msg("server stopped")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	cancel()
	log.Info().Msg("Shutdown Server ...")
	wg.Wait()
}
