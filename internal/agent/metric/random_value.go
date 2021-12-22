package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"math/rand"
)

type RandomValue struct {
	metric.GaugeBaseMetric
}

func (a *RandomValue) Name() string {
	return "RandomValue"
}

func (a *RandomValue) GetValue() interface{} {
	return metric.Gauge(rand.Float64())
}
