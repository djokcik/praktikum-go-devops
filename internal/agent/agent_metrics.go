package agent

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/agent/metric"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
	"net/http"
)

//go:generate mockery --name=AgentMetric

type agent struct {
	CollectedMetric map[string]SendAgentMetric
	Client          *http.Client
	metrics         []AgentMetric
	cfg             *Config
}

func NewAgent(cfg *Config) *agent {
	agent := new(agent)
	agent.CollectedMetric = make(map[string]SendAgentMetric)
	agent.Client = &http.Client{}
	agent.metrics = GetAgentMetrics()
	agent.cfg = cfg

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
		new(metric.TotalAlloc),
		new(metric.Mallocs),
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

func (a *agent) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "Agent").Logger()

	return &logger
}
