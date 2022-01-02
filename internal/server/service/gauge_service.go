package service

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=GaugeService

type GaugeService interface {
	GetOne(ctx context.Context, name string) (metric.Gauge, error)
	Update(ctx context.Context, name string, value metric.Gauge) (bool, error)
	List(ctx context.Context) ([]metric.Metric, error)
}

type GaugeServiceImpl struct {
	Repo storage.Repository
}

func (s *GaugeServiceImpl) Update(ctx context.Context, name string, value metric.Gauge) (bool, error) {
	val, err := s.Repo.Update(ctx, name, value)
	if err != nil {
		return val, err
	}

	s.Log(ctx).Info().Msg("metric updated")
	return val, nil
}

func (s *GaugeServiceImpl) GetOne(ctx context.Context, name string) (metric.Gauge, error) {
	val, err := s.Repo.Get(ctx, &storage.GetRepositoryFilter{
		Name: name,
		Type: metric.GaugeType,
	})

	if err != nil {
		return metric.Gauge(0), err
	}

	return val.(metric.Gauge), nil
}

func (s *GaugeServiceImpl) List(ctx context.Context) ([]metric.Metric, error) {
	metrics, err := s.Repo.List(ctx, &storage.ListRepositoryFilter{Type: metric.GaugeType})
	if err != nil {
		return nil, err
	}

	return metrics.([]metric.Metric), nil
}

func (s *GaugeServiceImpl) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "gauge metric service").Logger()

	return &logger
}
