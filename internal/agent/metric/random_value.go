package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"math/rand"
	"time"
)

type RandomValue struct {
	metric.GaugeBaseMetric
}

func (a *RandomValue) Name() string {
	return "RandomValue"
}

func (a *RandomValue) GetValue() interface{} {
	rand.Seed(time.Now().UnixNano())
	return metric.Gauge(rand.Float64())
}
