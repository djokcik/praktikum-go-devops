package service

import (
	"context"
	"errors"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=CounterService

type CounterService interface {
	GetOne(ctx context.Context, name string) (metric.Counter, error)
	Update(ctx context.Context, name string, value metric.Counter) (bool, error)
	List(ctx context.Context) ([]metric.Metric, error)
	Increase(ctx context.Context, name string, value metric.Counter) error
}

type CounterServiceImpl struct {
	Repo storage.Repository
}

func (s *CounterServiceImpl) GetOne(ctx context.Context, name string) (metric.Counter, error) {
	val, err := s.Repo.Get(ctx, &storage.GetRepositoryFilter{
		Name: name,
		Type: metric.CounterType,
	})

	if err != nil {
		return metric.Counter(0), err
	}

	return val.(metric.Counter), nil
}

func (s *CounterServiceImpl) Update(ctx context.Context, name string, value metric.Counter) (bool, error) {
	val, err := s.Repo.Update(ctx, name, value)
	if err != nil {
		return val, err
	}

	s.Log(ctx).Info().Msg("metric updated")
	return val, nil
}

func (s *CounterServiceImpl) List(ctx context.Context) ([]metric.Metric, error) {
	metrics, err := s.Repo.List(ctx, &storage.ListRepositoryFilter{Type: metric.CounterType})

	if err != nil {
		return nil, err
	}

	return metrics.([]metric.Metric), nil
}

func (s *CounterServiceImpl) Increase(ctx context.Context, name string, value metric.Counter) error {
	val, err := s.GetOne(ctx, name)

	if err != nil {
		if !errors.Is(err, storage.ErrValueNotFound) {
			return err
		}
	}

	_, err = s.Update(ctx, name, val+value)
	if err != nil {
		return errors.New("invalid save metric")
	}

	s.Log(ctx).Info().Msg("metric increased")

	return nil
}

func (s *CounterServiceImpl) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "counter metric service").Logger()

	return &logger
}
