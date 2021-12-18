package main

import (
	"github.com/Jokcik/praktikum-go-devops/internal/agent"
	"github.com/Jokcik/praktikum-go-devops/internal/agent/agentmetrics"
)

func main() {
	var updatedMetric = make(map[string]agentmetrics.AgentMetric)
	go agent.ReportMetricsToServer(updatedMetric)

	agent.CollectMetrics(updatedMetric)
}
