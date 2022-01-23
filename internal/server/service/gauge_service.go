package service

import (
	"context"
	"errors"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/djokcik/praktikum-go-devops/internal/service"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=GaugeService

type GaugeService interface {
	GetOne(ctx context.Context, name string) (metric.Gauge, error)
	Update(ctx context.Context, name string, value metric.Gauge) (bool, error)
	List(ctx context.Context) ([]metric.Metric, error)
	Verify(ctx context.Context, name string, value metric.Gauge, hash string) bool
}

type GaugeServiceImpl struct {
	Hash service.HashService
	Repo storage.MetricRepository
}

func (s GaugeServiceImpl) Update(ctx context.Context, name string, value metric.Gauge) (bool, error) {
	val, err := s.Repo.Update(ctx, name, value)
	if err != nil {
		return val, err
	}

	s.Log(ctx).Info().Msg("metric updated")
	return val, nil
}

func (s GaugeServiceImpl) GetOne(ctx context.Context, name string) (metric.Gauge, error) {
	val, err := s.Repo.Get(ctx, storage.GetRepositoryFilter{
		Name: name,
		Type: metric.GaugeType,
	})

	if err != nil {
		return metric.Gauge(0), err
	}

	if _, ok := val.(metric.Gauge); !ok {
		s.Log(ctx).Error().Msgf("value %v isn`t type Gauge", val)
		return metric.Gauge(0), errors.New("error parse gauge value")
	}

	return val.(metric.Gauge), nil
}

func (s GaugeServiceImpl) List(ctx context.Context) ([]metric.Metric, error) {
	metrics, err := s.Repo.List(ctx, storage.ListRepositoryFilter{Type: metric.GaugeType})
	if err != nil {
		return nil, err
	}

	return metrics.([]metric.Metric), nil
}

func (s GaugeServiceImpl) Verify(ctx context.Context, name string, value metric.Gauge, hash string) bool {
	actualHash := s.Hash.GetGaugeHash(ctx, name, value)
	return s.Hash.Verify(ctx, hash, actualHash)
}

func (s GaugeServiceImpl) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "gauge metric service").Logger()

	return &logger
}
