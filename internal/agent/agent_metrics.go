package agent

import (
	metric2 "github.com/Jokcik/praktikum-go-devops/internal/agent/metric"
	"net/http"
)

//go:generate mockery --name=AgentMetric

type agent struct {
	CollectedMetric map[string]SendAgentMetric
	Client          *http.Client
	metrics         []AgentMetric
}

func NewAgent() *agent {
	agent := new(agent)
	agent.CollectedMetric = make(map[string]SendAgentMetric)
	agent.Client = &http.Client{}
	agent.metrics = GetAgentMetrics()

	return agent
}

type SendAgentMetric struct {
	Name  string
	Type  string
	Value interface{}
}

type AgentMetric interface {
	Name() string
	Type() string
	GetValue() interface{}
}

func GetAgentMetrics() []AgentMetric {
	return []AgentMetric{
		// gauges
		new(metric2.Alloc),
		new(metric2.BuckHashSys),
		new(metric2.Frees),
		new(metric2.GCSys),
		new(metric2.GCCPUFraction),
		new(metric2.HeapAlloc),
		new(metric2.HeapIdle),
		new(metric2.HeapInuse),
		new(metric2.HeapObjects),
		new(metric2.HeapReleased),
		new(metric2.HeapSys),
		new(metric2.LastGC),
		new(metric2.Lookups),
		new(metric2.MCacheInuse),
		new(metric2.MCacheSys),
		new(metric2.MSpanInuse),
		new(metric2.MSpanSys),
		new(metric2.NextGC),
		new(metric2.NumForcedGC),
		new(metric2.NumGC),
		new(metric2.OtherSys),
		new(metric2.PauseTotalNs),
		new(metric2.StackInuse),
		new(metric2.StackSys),
		new(metric2.Sys),
		new(metric2.RandomValue),

		// Counter
		new(metric2.PollCount),
	}
}
