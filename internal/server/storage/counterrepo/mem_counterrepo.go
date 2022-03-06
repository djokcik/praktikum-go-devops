package counterrepo

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/storageconst"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/storer"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
	"sync"
)

type inmemDB struct {
	sync.RWMutex
	data map[string]metric.Counter
}

var (
	_ Repository = (*inmemRepository)(nil)
)

type inmemRepository struct {
	db    inmemDB
	store storer.MetricStorer
	cfg   server.Config
}

// NewInMem returns a new instance of Repository
func NewInMem(store storer.MetricStorer, cfg server.Config) Repository {
	data := make(map[string]metric.Counter)

	if store != nil {
		store.SetCounterDB(data)
	}

	return &inmemRepository{
		store: store,
		db:    inmemDB{data: data},
		cfg:   cfg,
	}
}

func (r *inmemRepository) notifyUpdateDBValue(ctx context.Context) {
	if r.store != nil && r.cfg.StoreInterval == 0 {
		r.store.SaveDBValue(ctx)
	}
}

func (r *inmemRepository) UpdateList(ctx context.Context, metrics []metric.CounterDto) error {
	r.db.Lock()
	defer r.db.Unlock()

	for _, counterMetric := range metrics {
		r.db.data[counterMetric.Name] = counterMetric.Value
	}

	r.notifyUpdateDBValue(ctx)
	r.Log(ctx).Info().Msg("metric updated")

	return nil
}

func (r *inmemRepository) Update(ctx context.Context, name string, value metric.Counter) error {
	r.db.Lock()
	defer r.db.Unlock()

	r.db.data[name] = value

	r.notifyUpdateDBValue(ctx)
	r.Log(ctx).Info().Msg("metric updated")

	return nil
}

func (r *inmemRepository) List(ctx context.Context) ([]metric.Metric, error) {
	var metricList []metric.Metric

	for metricName, metricValue := range r.db.data {
		metricList = append(metricList, metric.Metric{Name: metricName, Value: metricValue})
	}

	r.Log(ctx).Info().Msg("list finished")

	return metricList, nil
}

func (r *inmemRepository) Get(ctx context.Context, name string) (metric.Counter, error) {
	value, ok := r.db.data[name]
	if !ok {
		return 0, storageconst.ErrValueNotFound
	}

	r.Log(ctx).Info().Msg("get metric finished")

	return value, nil
}

func (r *inmemRepository) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "metric counter inmemRepository").Logger()

	return &logger
}
