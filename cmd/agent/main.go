package main

import (
	"github.com/Jokcik/praktikum-go-devops/internal/agent"
)

func main() {
	var updatedMetric = make(map[string]agent.SendAgentMetric)
	go agent.ReportMetricsToServer(updatedMetric)

	agent.CollectMetrics(updatedMetric)
}
