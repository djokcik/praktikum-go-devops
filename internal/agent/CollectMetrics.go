package agent

import (
	"time"
)

const pollInterval = 2

func CollectMetrics(updatedMetric map[string]SendAgentMetric) {
	var availableMetrics = GetAgentMetrics()
	ticker := time.NewTicker(pollInterval * time.Second)

	for range ticker.C {
		updateMetrics(updatedMetric, availableMetrics)
	}
}

func updateMetrics(updatedMetric map[string]SendAgentMetric, availableMetrics []AgentMetric) {
	for _, metric := range availableMetrics {
		var name = metric.Name()
		updatedMetric[name] = SendAgentMetric{Metric: metric, Value: metric.GetValue()}
	}
}
