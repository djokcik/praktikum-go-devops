package agentmetrics

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metrics"
)

type SendAgentMetric struct {
	Metric metrics.Metric
	Value  interface{}
}

type AgentMetric interface {
	Type() string
	Name() string
	GetValue() interface{}
}

func GetAgentMetrics() []AgentMetric {
	return []AgentMetric{
		// gauges
		new(metrics.Alloc),
		new(metrics.BuckHashSys),
		new(metrics.Frees),
		new(metrics.GCSys),
		new(metrics.GCCPUFraction),
		new(metrics.HeapAlloc),
		new(metrics.HeapIdle),
		new(metrics.HeapInuse),
		new(metrics.HeapObjects),
		new(metrics.HeapReleased),
		new(metrics.HeapSys),
		new(metrics.LastGC),
		new(metrics.Lookups),
		new(metrics.MCacheInuse),
		new(metrics.MCacheSys),
		new(metrics.MSpanInuse),
		new(metrics.MSpanSys),
		new(metrics.NextGC),
		new(metrics.NumForcedGC),
		new(metrics.NumGC),
		new(metrics.OtherSys),
		new(metrics.PauseTotalNs),
		new(metrics.StackInuse),
		new(metrics.StackSys),
		new(metrics.Sys),

		// Counter
		new(metrics.PollCount),
	}
}
