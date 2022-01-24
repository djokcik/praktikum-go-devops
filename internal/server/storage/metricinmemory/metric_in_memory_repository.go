package metricinmemory

import (
	"context"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/helpers"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/metricinmemory/storer"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/metricinmemory/storer/filestorer"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/model"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
	"sync"
)

type repository struct {
	store      storer.MetricStorer
	cfg        server.Config
	inMemoryDB *model.InMemoryMetricDB
}

func NewRepository(ctx context.Context, wg *sync.WaitGroup, cfg server.Config) storage.MetricRepository {
	r := &repository{}

	r.cfg = cfg
	r.inMemoryDB = new(model.InMemoryMetricDB)
	r.inMemoryDB.CounterMapMetric = make(map[string]metric.Counter)
	r.inMemoryDB.GaugeMapMetric = make(map[string]metric.Gauge)

	if cfg.StoreFile != "" {
		r.store = filestorer.NewMetricFileStorer(ctx, r.inMemoryDB, cfg)
	} else {
		r.Log(ctx).Info().Msg("saving metrics to the store is disabled")
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

func (r *repository) notifyUpdateDBValue(ctx context.Context) {
	if r.store != nil && r.cfg.StoreInterval == 0 {
		r.store.SaveDBValue(ctx)
	}
}

func (r *repository) Update(ctx context.Context, name string, value interface{}) (bool, error) {
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

func (r repository) List(ctx context.Context, filter storage.ListRepositoryFilter) (interface{}, error) {
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

func (r repository) Get(ctx context.Context, filter storage.GetRepositoryFilter) (interface{}, error) {
	metricType := filter.Type
	var value interface{}
	var ok bool

	switch metricType {
	default:
		return nil, fmt.Errorf("type `%v` isn`t avalilable metric type", filter.Type)
	case metric.GaugeType:
		value, ok = r.inMemoryDB.GaugeMapMetric[filter.Name]
		if !ok {
			return 0, storage.ErrValueNotFound
		}
	case metric.CounterType:
		value, ok = r.inMemoryDB.CounterMapMetric[filter.Name]
		if !ok {
			return 0, storage.ErrValueNotFound
		}
	}

	r.Log(ctx).Info().Msg("get metric finished")

	return value, nil
}

func (r repository) Ping(ctx context.Context) error {
	return nil
}

func (r *repository) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "metric repository").Logger()

	return &logger
}
