package agent

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metrics"
	"time"
)

const pollInterval = 2

func CollectMetrics(updatedMetric *map[string]metrics.Metric) {
	var availableMetrics = metrics.GetMetrics()

	for {
		for _, metric := range availableMetrics {
			var name = metric.Name()
			(*updatedMetric)[name] = metric
		}

		time.Sleep(pollInterval * time.Second)
	}
}
