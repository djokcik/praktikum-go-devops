package servermetrics

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metrics"
)

type ServerMapMetrics map[string]metrics.Metric

func GetServerGaugeMetrics() ServerMapMetrics {
	mapMetrics := make(ServerMapMetrics)

	var gaugeMetrics = metrics.GetGaugesMetrics()
	for _, metric := range gaugeMetrics {
		mapMetrics[metric.Name()] = metric
	}

	return mapMetrics
}

func GetServerCounterMetrics() ServerMapMetrics {
	mapMetrics := make(ServerMapMetrics)

	var counterMetrics = metrics.GetCounterMetrics()
	for _, metric := range counterMetrics {
		mapMetrics[metric.Name()] = metric
	}

	return mapMetrics
}
