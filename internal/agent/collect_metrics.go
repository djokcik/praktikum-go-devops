package agent

import (
	"github.com/Jokcik/praktikum-go-devops/internal/agent/agentmetrics"
	"time"
)

const pollInterval = 2

func CollectMetrics(updatedMetric map[string]agentmetrics.AgentMetric) {
	var availableMetrics = agentmetrics.GetAgentMetrics()
	ticker := time.NewTicker(pollInterval * time.Second)

	for {
		<-ticker.C

		for _, metric := range availableMetrics {
			var name = metric.Name()
			updatedMetric[name] = metric
		}
	}
}
