package main

import (
	"context"
	appgrpc "github.com/djokcik/praktikum-go-devops/internal/grpc"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/handler"
	"github.com/djokcik/praktikum-go-devops/internal/server/service"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/reporegistry"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	pb "github.com/djokcik/praktikum-go-devops/pkg/proto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"sync"
)

func getRepoRegistry(ctx context.Context, cfg server.Config, wg *sync.WaitGroup) reporegistry.RepoRegistry {
	var repoRegistry reporegistry.RepoRegistry

	if cfg.DatabaseDsn != "" {
		databaseService := service.NewDatabaseService(ctx, cfg)
		db, err := databaseService.Open(ctx, wg)
		if err != nil {
			logging.NewLogger().Fatal().Err(err).Msgf("Doesn`t open database connection")
			os.Exit(1)
		}

		repoRegistry = reporegistry.NewPostgreSQL(ctx, db)
	} else {
		repoRegistry = reporegistry.NewInMem(ctx, wg, cfg)
	}

	return repoRegistry
}

func makeGRPCMetricServer(cfg server.Config, registry reporegistry.RepoRegistry) (*grpc.Server, net.Listener, error) {
	listen, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		logging.NewLogger().Fatal().Err(err).Msgf("error listen gRPC server with address: %s", err)
		return nil, nil, err
	}

	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterMetricsServer(s, appgrpc.MakeGRPCMetricService(cfg, registry))

	return s, listen, nil
}

func makeMetricRoutes(_ context.Context, repoRegistry reporegistry.RepoRegistry, mux *chi.Mux, cfg server.Config) *handler.Handler {
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

	h.Mount("/debug", middleware.Profiler())

	return h
}
