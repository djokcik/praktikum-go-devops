package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/helpers"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/model"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/store"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
	"sync"
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
	Ping(ctx context.Context) error
}

var (
	ErrValueNotFound = errors.New("value not found")
)

type MetricRepositoryImpl struct {
	store      store.MetricStorer
	inMemoryDB *model.InMemoryMetricDB
	cfg        server.Config
}

func NewMetricRepository(ctx context.Context, wg *sync.WaitGroup, cfg server.Config) MetricRepository {
	r := &MetricRepositoryImpl{}

	r.cfg = cfg
	r.inMemoryDB = new(model.InMemoryMetricDB)
	r.inMemoryDB.CounterMapMetric = make(map[string]metric.Counter)
	r.inMemoryDB.GaugeMapMetric = make(map[string]metric.Gauge)

	if cfg.DatabaseDsn != "" {
		r.store = store.NewMetricDBStorer(ctx, r.inMemoryDB, cfg)
	} else if cfg.StoreFile != "" {
		r.store = store.NewMetricFileStorer(ctx, r.inMemoryDB, cfg)
	} else {
		r.Log(ctx).Info().Msg("save metrics to store are disabled")
		r.store = nil
	}

	if r.store != nil {
		wg.Add(1)
		go func() {
			<-ctx.Done()
			defer wg.Done()

			r.store.SaveDBValue(ctx)
			r.store.Close()
		}()

		if cfg.Restore {
			r.store.RestoreDBValue(ctx)
		}

		if cfg.StoreInterval != 0 {
			go helpers.SetTicker(func() { r.store.SaveDBValue(ctx) }, cfg.StoreInterval)
		}
	}

	return r
}

func (r *MetricRepositoryImpl) notifyUpdateDBValue(ctx context.Context) {
	if r.store != nil && r.cfg.StoreInterval == 0 {
		r.store.SaveDBValue(ctx)
	}
}

func (r *MetricRepositoryImpl) Update(ctx context.Context, name string, value interface{}) (bool, error) {
	r.inMemoryDB.Lock()
	defer r.inMemoryDB.Unlock()

	switch metricValue := value.(type) {
	default:
		return false, fmt.Errorf("entity could`t convert `%v` into available metric type", value)
	case metric.Counter:
		r.inMemoryDB.CounterMapMetric[name] = metricValue
	case metric.Gauge:
		r.inMemoryDB.GaugeMapMetric[name] = metricValue
	}

	r.notifyUpdateDBValue(ctx)
	r.Log(ctx).Info().Msg("metric updated")

	return true, nil
}

func (r MetricRepositoryImpl) List(ctx context.Context, filter ListRepositoryFilter) (interface{}, error) {
	var metricList []metric.Metric

	switch filter.Type {
	default:
		return nil, fmt.Errorf("type `%v` isn`t avalilable metric type", filter.Type)
	case metric.GaugeType:
		for metricName, metricValue := range r.inMemoryDB.GaugeMapMetric {
			metricList = append(metricList, metric.Metric{Name: metricName, Value: metricValue})
		}
	case metric.CounterType:
		for metricName, metricValue := range r.inMemoryDB.CounterMapMetric {
			metricList = append(metricList, metric.Metric{Name: metricName, Value: metricValue})
		}
	}

	r.Log(ctx).Info().Msg("list finished")

	return metricList, nil
}

func (r MetricRepositoryImpl) Ping(ctx context.Context) error {
	return r.store.Ping(ctx)
}

func (r MetricRepositoryImpl) Get(ctx context.Context, filter GetRepositoryFilter) (interface{}, error) {
	metricType := filter.Type
	var value interface{}
	var ok bool

	switch metricType {
	default:
		return nil, fmt.Errorf("type `%v` isn`t avalilable metric type", filter.Type)
	case metric.GaugeType:
		value, ok = r.inMemoryDB.GaugeMapMetric[filter.Name]
		if !ok {
			return 0, ErrValueNotFound
		}
	case metric.CounterType:
		value, ok = r.inMemoryDB.CounterMapMetric[filter.Name]
		if !ok {
			return 0, ErrValueNotFound
		}
	}

	r.Log(ctx).Info().Msg("get metric finished")

	return value, nil
}

func (r *MetricRepositoryImpl) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "in memory metric repository").Logger()

	return &logger
}
