package storer

import "context"

//go:generate mockery --name=MetricStorer

type MetricStorer interface {
	RestoreDBValue(ctx context.Context)
	SaveDBValue(ctx context.Context)
	Close()
}
