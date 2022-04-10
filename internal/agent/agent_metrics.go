package agent

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/agent/metric"
	"github.com/djokcik/praktikum-go-devops/internal/service"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	pb "github.com/djokcik/praktikum-go-devops/pkg/proto"
	"github.com/rs/zerolog"
	"github.com/shirou/gopsutil/cpu"
	"google.golang.org/grpc"
	"net/http"
	"sync"
)

//go:generate mockery --name=AgentMetric

var (
	_ Agent = (*agent)(nil)
)

type Agent interface {
	CollectMetrics(ctx context.Context) func()
	CollectPsutilMetrics(ctx context.Context) func()
	SendToServer(ctx context.Context)
}

type agent struct {
	sync.RWMutex
	CollectedMetric map[string]SendAgentMetric
	Client          *http.Client
	Hash            service.HashService
	metrics         []AgentMetric
	cfg             Config
	GRPCClient      pb.MetricsClient
}

func NewAgent(cfg Config, conn grpc.ClientConnInterface) Agent {
	metricAgent := new(agent)
	metricAgent.CollectedMetric = make(map[string]SendAgentMetric)
	metricAgent.Client = &http.Client{}
	metricAgent.metrics = GetAgentMetrics()
	metricAgent.cfg = cfg
	metricAgent.Hash = service.NewHashService(cfg.Key)

	if conn != nil {
		metricAgent.GRPCClient = pb.NewMetricsClient(conn)
	}

	return metricAgent
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

func GetAgentPsutilMetrics() []AgentMetric {
	var agents []AgentMetric
	percent, err := cpu.Percent(0, true)
	if err == nil {
		for i := 0; i < len(percent); i++ {
			agents = append(agents, metric.NewCPUUtilization(i, percent[i]))
		}
	}

	return append(agents, new(metric.FreeMemory), new(metric.TotalMemory))
}

func (a *agent) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "Agent").Logger()

	return &logger
}
