package metric

import (
	"math/rand"
)

type RandomValue struct {
	GaugeBaseMetric
}

func (a *RandomValue) Name() string {
	return "RandomValue"
}

func (a *RandomValue) GetValue() interface{} {
	return Gauge(rand.Float64())
}
