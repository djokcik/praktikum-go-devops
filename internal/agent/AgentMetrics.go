package agent

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
)

//go:generate mockery --name=AgentMetric

type SendAgentMetric struct {
	Metric metric.Metric
	Value  interface{}
}

type AgentMetric interface {
	metric.Metric
	GetValue() interface{}
}

func GetAgentMetrics() []AgentMetric {
	return []AgentMetric{
		// gauges
		new(metric.Alloc),
		new(metric.BuckHashSys),
		new(metric.Frees),
		new(metric.GCSys),
		new(metric.GCCPUFraction),
		new(metric.HeapAlloc),
		new(metric.HeapIdle),
		new(metric.HeapInuse),
		new(metric.HeapObjects),
		new(metric.HeapReleased),
		new(metric.HeapSys),
		new(metric.LastGC),
		new(metric.Lookups),
		new(metric.MCacheInuse),
		new(metric.MCacheSys),
		new(metric.MSpanInuse),
		new(metric.MSpanSys),
		new(metric.NextGC),
		new(metric.NumForcedGC),
		new(metric.NumGC),
		new(metric.OtherSys),
		new(metric.PauseTotalNs),
		new(metric.StackInuse),
		new(metric.StackSys),
		new(metric.Sys),
		new(metric.RandomValue),

		// Counter
		new(metric.PollCount),
	}
}
