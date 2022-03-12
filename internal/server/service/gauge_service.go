package service

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/gaugerepo"
	"github.com/djokcik/praktikum-go-devops/internal/service"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=GaugeService

// GaugeService provide methods for control gauge metrics
type GaugeService interface {
	GetOne(ctx context.Context, name string) (metric.Gauge, error)
	Update(ctx context.Context, name string, value metric.Gauge) error
	UpdateList(ctx context.Context, metrics []metric.GaugeDto) error
	List(ctx context.Context) ([]metric.Metric, error)
	Verify(ctx context.Context, name string, value metric.Gauge, hash string) bool
}

var (
	_ GaugeService = (*GaugeServiceImpl)(nil)
)

type GaugeServiceImpl struct {
	Hash service.HashService
	Repo gaugerepo.Repository
}

// Update - update gauge metric with name and value
func (s GaugeServiceImpl) Update(ctx context.Context, name string, value metric.Gauge) error {
	err := s.Repo.Update(ctx, name, value)
	if err != nil {
		return err
	}

	s.Log(ctx).Info().Msg("metric updated")
	return nil
}

// UpdateList - update list gauge metrics with []metric.GaugeDto
func (s GaugeServiceImpl) UpdateList(ctx context.Context, metrics []metric.GaugeDto) error {
	err := s.Repo.UpdateList(ctx, metrics)

	if err != nil {
		return err
	}

	return nil
}

// GetOne - return metric by name
func (s GaugeServiceImpl) GetOne(ctx context.Context, name string) (metric.Gauge, error) {
	val, err := s.Repo.Get(ctx, name)

	if err != nil {
		return metric.Gauge(0), err
	}

	return val, nil
}

// List - return list of gauge metrics
func (s GaugeServiceImpl) List(ctx context.Context) ([]metric.Metric, error) {
	metrics, err := s.Repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// Verify - check equal hash and calculate metric hash
func (s GaugeServiceImpl) Verify(ctx context.Context, name string, value metric.Gauge, hash string) bool {
	actualHash := s.Hash.GetGaugeHash(ctx, name, value)
	return s.Hash.Verify(ctx, hash, actualHash)
}

func (s GaugeServiceImpl) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "gauge metric service").Logger()

	return &logger
}
