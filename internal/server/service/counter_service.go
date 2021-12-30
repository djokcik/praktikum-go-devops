package service

import (
	"errors"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
)

//go:generate mockery --name=CounterService

type CounterService interface {
	GetOne(name string) (metric.Counter, error)
	Update(name string, value metric.Counter) (bool, error)
	List() ([]metric.Metric, error)
	AddValue(name string, value metric.Counter) error
}

type CounterServiceImpl struct {
	Repo storage.Repository
}

func (s *CounterServiceImpl) GetOne(name string) (metric.Counter, error) {
	val, err := s.Repo.Get(&storage.GetRepositoryFilter{
		Name: name,
		Type: metric.CounterType,
	})

	if err != nil {
		return metric.Counter(0), err
	}

	return val.(metric.Counter), nil
}

func (s *CounterServiceImpl) Update(name string, value metric.Counter) (bool, error) {
	return s.Repo.Update(name, value)
}

func (s *CounterServiceImpl) List() ([]metric.Metric, error) {
	metrics, err := s.Repo.List(&storage.ListRepositoryFilter{Type: metric.CounterType})

	if err != nil {
		return nil, err
	}

	return metrics.([]metric.Metric), nil
}

func (s *CounterServiceImpl) AddValue(name string, value metric.Counter) error {
	val, err := s.GetOne(name)

	if err != nil {
		if err.Error() != storage.ValueNotFound {
			return err
		}
	}

	_, err = s.Update(name, val+value)
	if err != nil {
		return errors.New("invalid save metric")
	}

	return nil
}
