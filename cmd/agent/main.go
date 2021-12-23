package main

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/agent"
	"github.com/djokcik/praktikum-go-devops/internal/helpers"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const reportInterval = 10 * time.Second
const pollInterval = 2 * time.Second

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metricAgent := agent.NewAgent(ctx)

	go helpers.SetTicker(metricAgent.CollectMetrics, pollInterval)
	go helpers.SetTicker(metricAgent.SendToServer, reportInterval)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Agent ...")
}
