package main

import (
	"github.com/Jokcik/praktikum-go-devops/internal/agent"
	metrics "github.com/Jokcik/praktikum-go-devops/internal/metrics"
)

func main() {
	var updatedMetric = make(map[string]metrics.Metric)
	go agent.ReportMetricsToServer(&updatedMetric)

	agent.CollectMetrics(&updatedMetric)
}
