package main

import (
	"github.com/Jokcik/praktikum-go-devops/internal/agent"
	"github.com/Jokcik/praktikum-go-devops/internal/helpers"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const reportInterval = 10
const pollInterval = 2

func main() {
	metricAgent := agent.NewAgent()

	go helpers.SetTicker(metricAgent.CollectMetrics, pollInterval)
	go helpers.SetTicker(metricAgent.SendToServer, reportInterval)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Agent ...")
}
