package main

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/agent"
	"github.com/djokcik/praktikum-go-devops/internal/helpers"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"google.golang.org/grpc"
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
	var wg sync.WaitGroup

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	var cfg = agent.NewConfig()
	log := logging.NewLogger()

	log.Info().Msgf("Build version: %s", buildVersion)
	log.Info().Msgf("Build date: %s", buildDate)
	log.Info().Msgf("Build commit: %s", buildCommit)

	logging.
		NewLogger().
		Info().
		Msgf("Config: %+v", cfg)

	var conn grpc.ClientConnInterface

	if cfg.GRPCAddress != "" {
		var err error

		conn, err = grpc.Dial(cfg.GRPCAddress)
		if err != nil {
			logging.NewLogger().Fatal().Err(err).Msg("invalid connect to gRPC server")
		}
	}

	metricAgent := agent.NewAgent(cfg, conn)

	go helpers.SetTicker(metricAgent.CollectPsutilMetrics(ctx), cfg.PollPsutilsInterval)
	go helpers.SetTicker(metricAgent.CollectMetrics(ctx), cfg.PollInterval)
	go helpers.SetTicker(func() {
		wg.Add(1)
		defer wg.Done()

		metricAgent.SendToServer(ctx)
	}, cfg.ReportInterval)

	<-ctx.Done()
	log.Info().Msg("Shutdown Agent ...")
	wg.Wait()
}
