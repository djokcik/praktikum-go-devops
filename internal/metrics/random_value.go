package metrics

import (
	"math/rand"
)

type randomValue struct {
	GaugeBaseMetric
}

func (a randomValue) Name() string {
	return "RandomValue"
}

func (a randomValue) GetValue() interface{} {
	return Gauge(rand.Float64())
}
