package grpc

import (
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/service"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/reporegistry"
	commonService "github.com/djokcik/praktikum-go-devops/internal/service"
)

func MakeGRPCMetricService(cfg server.Config, registry reporegistry.RepoRegistry) *MetricsServer {
	hashService := commonService.NewHashService(cfg.Key)

	return &MetricsServer{
		Counter: &service.CounterServiceImpl{Repo: registry.GetCounterRepo(), Hash: hashService},
		Gauge:   &service.GaugeServiceImpl{Repo: registry.GetGaugeRepo(), Hash: hashService},
	}
}
