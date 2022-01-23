package model

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"sync"
)

type InMemoryMetricDB struct {
	sync.RWMutex
	CounterMapMetric map[string]metric.Counter
	GaugeMapMetric   map[string]metric.Gauge
}
