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

type MetricInMemoryRepository struct {
	Store store.MetricStore
	db    *model.InMemoryDatabase
}

func NewMetricInMemoryRepository(ctx context.Context, wg *sync.WaitGroup, cfg server.Config) MetricRepository {
	r := &MetricInMemoryRepository{}
	r.db = new(model.InMemoryDatabase)
	r.db.CounterMapMetric = make(map[string]metric.Counter)
	r.db.GaugeMapMetric = make(map[string]metric.Gauge)

	r.Store = &store.MetricStoreFile{DB: r.db, Cfg: cfg}
	r.Store.Configure(ctx, wg)

	return r
}

func (r *MetricInMemoryRepository) Update(ctx context.Context, name string, value interface{}) (bool, error) {
	r.db.Lock()
	defer r.db.Unlock()

	switch metricValue := value.(type) {
	default:
		return false, fmt.Errorf("entity could`t convert `%v` into available metric type", value)
	case metric.Counter:
		r.db.CounterMapMetric[name] = metricValue
	case metric.Gauge:
		r.db.GaugeMapMetric[name] = metricValue
	}

	r.Store.NotifyUpdateDBValue(ctx)
	r.Log(ctx).Info().Msg("metric updated")

	return true, nil
}

func (r MetricInMemoryRepository) List(ctx context.Context, filter ListRepositoryFilter) (interface{}, error) {
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

func (r MetricInMemoryRepository) Get(ctx context.Context, filter GetRepositoryFilter) (interface{}, error) {
	metricType := filter.Type
	var value interface{}
	var ok bool

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

func (r *MetricInMemoryRepository) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "metric repository").Logger()

	return &logger
}
