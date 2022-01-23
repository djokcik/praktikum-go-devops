package storage

import (
	"context"
)

//go:generate mockery --name=MetricRepository

type GetRepositoryFilter struct {
	Type string
	Name string
}

type ListRepositoryFilter struct {
	Type string
}

type MetricRepository interface {
	Update(ctx context.Context, name string, entity interface{}) (bool, error)
	List(ctx context.Context, filter ListRepositoryFilter) (interface{}, error)
	Get(ctx context.Context, filter GetRepositoryFilter) (interface{}, error)
}
