package storer

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
)

//go:generate mockery --name=MetricStorer

type MetricStorer interface {
	RestoreDBValue(ctx context.Context)
	SaveDBValue(ctx context.Context)
	SetCounterDB(map[string]metric.Counter)
	SetGaugeDB(map[string]metric.Gauge)
	Close()
}
