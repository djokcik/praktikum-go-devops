package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/djokcik/praktikum-go-devops/internal/agent"
	"github.com/djokcik/praktikum-go-devops/internal/helpers"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cfg agent.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	metricAgent := agent.NewAgent(&cfg)

	go helpers.SetTicker(metricAgent.CollectMetrics(ctx), cfg.PollInterval)
	go helpers.SetTicker(metricAgent.SendToServer(ctx), cfg.ReportInterval)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Agent ...")
}
