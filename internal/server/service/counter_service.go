package service

import (
	"context"
	"errors"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/counterrepo"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/storageconst"
	"github.com/djokcik/praktikum-go-devops/internal/service"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=CounterService

type CounterService interface {
	GetOne(ctx context.Context, name string) (metric.Counter, error)
	Update(ctx context.Context, name string, value metric.Counter) error
	UpdateList(ctx context.Context, metrics []metric.CounterDto) error
	List(ctx context.Context) ([]metric.Metric, error)
	Increase(ctx context.Context, name string, value metric.Counter) error
	Verify(ctx context.Context, name string, value metric.Counter, hash string) bool
}

type CounterServiceImpl struct {
	Hash service.HashService
	Repo counterrepo.Repository
}

func (s CounterServiceImpl) GetOne(ctx context.Context, name string) (metric.Counter, error) {
	val, err := s.Repo.Get(ctx, name)

	if err != nil {
		return metric.Counter(0), err
	}

	return val, nil
}

func (s CounterServiceImpl) Update(ctx context.Context, name string, value metric.Counter) error {
	err := s.Repo.Update(ctx, name, value)
	if err != nil {
		return err
	}

	s.Log(ctx).Info().Msg("metric updated")
	return nil
}

func (s CounterServiceImpl) UpdateList(ctx context.Context, metrics []metric.CounterDto) error {
	counterMap := make(map[string]metric.Counter)

	for _, dto := range metrics {
		counterMap[dto.Name] += dto.Value
	}

	metrics = make([]metric.CounterDto, 0)
	for name, value := range counterMap {
		metrics = append(metrics, metric.CounterDto{Name: name, Value: value})
	}

	err := s.Repo.UpdateList(ctx, metrics)

	if err != nil {
		return err
	}

	return nil
}

func (s CounterServiceImpl) List(ctx context.Context) ([]metric.Metric, error) {
	metrics, err := s.Repo.List(ctx)

	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (s CounterServiceImpl) Increase(ctx context.Context, name string, value metric.Counter) error {
	val, err := s.GetOne(ctx, name)

	if err != nil {
		if !errors.Is(err, storageconst.ErrValueNotFound) {
			return err
		}
	}

	err = s.Update(ctx, name, val+value)
	if err != nil {
		return errors.New("invalid save metric")
	}

	s.Log(ctx).Info().Msg("metric increased")

	return nil
}

func (s CounterServiceImpl) Verify(ctx context.Context, name string, value metric.Counter, hash string) bool {
	actualHash := s.Hash.GetCounterHash(ctx, name, value)
	return s.Hash.Verify(ctx, hash, actualHash)
}

func (s CounterServiceImpl) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "counter metric service").Logger()

	return &logger
}
