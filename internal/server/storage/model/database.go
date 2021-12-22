package model

import "github.com/Jokcik/praktikum-go-devops/internal/metric"

type Database struct {
	CounterMapMetric map[string]metric.Counter
	GaugeMapMetric   map[string]metric.Gauge
}
