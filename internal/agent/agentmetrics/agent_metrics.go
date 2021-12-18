package agentmetrics

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metrics"
)

type AgentMetric interface {
	Type() string
	Name() string
	GetValue() interface{}
}

func GetAgentMetrics() []AgentMetric {
	return []AgentMetric{
		// gauges
		metrics.BuckHashSys{},
		metrics.Frees{},
		metrics.GCSys{},
		metrics.GCCPUFraction{},
		metrics.HeapAlloc{},
		metrics.HeapIdle{},
		metrics.HeapInuse{},
		metrics.HeapObjects{},
		metrics.HeapReleased{},
		metrics.HeapSys{},
		metrics.LastGC{},
		metrics.Lookups{},
		metrics.MCacheInuse{},
		metrics.MCacheSys{},
		metrics.MSpanInuse{},
		metrics.MSpanSys{},
		metrics.NextGC{},
		metrics.NumForcedGC{},
		metrics.NumGC{},
		metrics.OtherSys{},
		metrics.PauseTotalNs{},
		metrics.StackInuse{},
		metrics.StackSys{},
		metrics.Sys{},

		// Counter
		metrics.PollCount{},
	}
}
