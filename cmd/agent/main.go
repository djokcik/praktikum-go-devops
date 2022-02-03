package main

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/agent"
	"github.com/djokcik/praktikum-go-devops/internal/helpers"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cfg = agent.NewConfig()

	logging.
		NewLogger().
		Info().
		Msgf("Config: %+v", cfg)

	metricAgent := agent.NewAgent(cfg)

	//go helpers.SetTicker(metricAgent.CollectPsutilMetrics(ctx), cfg.PollInterval)
	go helpers.SetTicker(metricAgent.CollectMetrics(ctx), cfg.PollInterval)
	go helpers.SetTicker(metricAgent.SendToServer(ctx), cfg.ReportInterval)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	logging.NewLogger().Info().Msg("Shutdown Agent ...")
}
