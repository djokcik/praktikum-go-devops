package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/model"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/store"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
	"sync"
)

var (
	ErrValueNotFound = errors.New("value not found")
)

type MetricRepository struct {
	BaseRepository
	Store store.MetricStore
}

func (r *MetricRepository) Configure(ctx context.Context, wg *sync.WaitGroup, db *model.Database, cfg server.Config) {
	r.BaseRepository.Configure(ctx, wg, db, cfg)

	r.Store = &store.MetricStoreFile{DB: db, Cfg: cfg}
	r.Store.Configure(ctx, wg)
}

func (r *MetricRepository) Update(ctx context.Context, name interface{}, value interface{}) (bool, error) {
	r.db.Lock()
	defer r.db.Unlock()

	switch metricValue := value.(type) {
	default:
		return false, fmt.Errorf("entity could`t convert `%v` into available metric type", value)
	case metric.Counter:
		r.db.CounterMapMetric[name.(string)] = metricValue
	case metric.Gauge:
		r.db.GaugeMapMetric[name.(string)] = metricValue
	}

	r.Store.NotifyUpdateDBValue(ctx)
	r.Log(ctx).Info().Msg("metric updated")

	return true, nil
}

func (r *MetricRepository) List(ctx context.Context, filter *ListRepositoryFilter) (interface{}, error) {
	var metricList []metric.Metric

	switch filter.Type {
	default:
		return nil, fmt.Errorf("type `%v` isn`t avalilable metric type", filter.Type)
	case metric.GaugeType:
		for metricName, metricValue := range r.db.GaugeMapMetric {
			metricList = append(metricList, metric.Metric{Name: metricName, Value: metricValue})
		}
	case metric.CounterType:
		for metricName, metricValue := range r.db.CounterMapMetric {
			metricList = append(metricList, metric.Metric{Name: metricName, Value: metricValue})
		}
	}

	r.Log(ctx).Info().Msg("list finished")

	return metricList, nil
}

func (r *MetricRepository) Get(ctx context.Context, filter *GetRepositoryFilter) (interface{}, error) {
	metricType := filter.Type
	var value interface{}
	var ok bool

	// TODO use db that code will be more simply
	switch metricType {
	default:
		return nil, fmt.Errorf("type `%v` isn`t avalilable metric type", filter.Type)
	case metric.GaugeType:
		value, ok = r.db.GaugeMapMetric[filter.Name]
		if !ok {
			return 0, ErrValueNotFound
		}
	case metric.CounterType:
		value, ok = r.db.CounterMapMetric[filter.Name]
		if !ok {
			return 0, ErrValueNotFound
		}
	}

	r.Log(ctx).Info().Msg("get metric finished")

	return value, nil
}

func (r *MetricRepository) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "metric repository").Logger()

	return &logger
}
