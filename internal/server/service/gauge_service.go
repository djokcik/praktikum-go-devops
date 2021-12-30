package service

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
)

//go:generate mockery --name=GaugeService

type GaugeService interface {
	GetOne(name string) (metric.Gauge, error)
	Update(name string, value metric.Gauge) (bool, error)
	List() ([]metric.Metric, error)
}

type GaugeServiceImpl struct {
	Repo storage.Repository
}

func (s *GaugeServiceImpl) Update(name string, value metric.Gauge) (bool, error) {
	return s.Repo.Update(name, value)
}

func (s *GaugeServiceImpl) GetOne(name string) (metric.Gauge, error) {
	val, err := s.Repo.Get(&storage.GetRepositoryFilter{
		Name: name,
		Type: metric.GaugeType,
	})

	if err != nil {
		return metric.Gauge(0), err
	}

	return val.(metric.Gauge), nil
}

func (s *GaugeServiceImpl) List() ([]metric.Metric, error) {
	metrics, err := s.Repo.List(&storage.ListRepositoryFilter{Type: metric.GaugeType})
	if err != nil {
		return nil, err
	}

	return metrics.([]metric.Metric), nil
}
