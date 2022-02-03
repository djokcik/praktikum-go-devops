package counterrepo

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
)

//go:generate mockery --name=Repository

type Repository interface {
	Get(ctx context.Context, name string) (metric.Counter, error)
	List(ctx context.Context) ([]metric.Metric, error)
	Update(ctx context.Context, name string, value metric.Counter) error
	UpdateList(ctx context.Context, metrics []metric.CounterDto) error
}
