package metricdatabase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/djokcik/praktikum-go-devops/migration"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
)

type repository struct {
	db      *sql.DB
	counter counterDatabaseService
	gauge   gaugeDatabaseService
}

func NewRepository(ctx context.Context, cfg server.Config) (storage.MetricRepository, error) {
	r := &repository{}

	db, err := sql.Open("pgx", cfg.DatabaseDsn)
	if err != nil {
		r.Log(ctx).Fatal().Err(err).Msgf("Unable to connect to database")
		return nil, err
	}

	r.db = db
	r.counter = &counterDatabaseServiceImpl{db: db}
	r.gauge = &gaugeDatabaseServiceImpl{db: db}

	err = migration.CreateCounterTable(db)
	if err != nil {
		r.Log(ctx).Warn().Err(err).Msgf("couldn't create counter table")
	}

	err = migration.CreateGaugeTable(db)
	if err != nil {
		r.Log(ctx).Warn().Err(err).Msgf("couldn't create gauge table")
	}

	go func() {
		<-ctx.Done()
		r.db.Close()
	}()

	return r, nil
}

func (r repository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r repository) Update(ctx context.Context, name string, value interface{}) (bool, error) {
	switch metricValue := value.(type) {
	default:
		return false, fmt.Errorf("entity could`t convert `%v` into available metric type", value)
	case metric.Counter:
		err := r.counter.Update(ctx, name, metricValue)
		if err != nil {
			return false, err
		}

	case metric.Gauge:
		err := r.gauge.Update(ctx, name, metricValue)
		if err != nil {
			return false, err
		}

	}

	r.Log(ctx).Info().Msg("metric updated")

	return true, nil
}

func (r repository) List(ctx context.Context, filter storage.ListRepositoryFilter) (interface{}, error) {
	var err error
	var metricList []metric.Metric

	switch filter.Type {
	default:
		return nil, fmt.Errorf("type `%v` isn`t avalilable metric type", filter.Type)
	case metric.CounterType:
		metricList, err = r.counter.List(ctx)

		if err != nil {
			return nil, err
		}

	case metric.GaugeType:
		metricList, err = r.gauge.List(ctx)

		if err != nil {
			return nil, err
		}
	}

	r.Log(ctx).Info().Msg("list finished")

	return metricList, nil
}

func (r repository) Get(ctx context.Context, filter storage.GetRepositoryFilter) (interface{}, error) {
	metricType := filter.Type
	var value interface{}
	var err error

	switch metricType {
	default:
		return nil, fmt.Errorf("type `%v` isn`t avalilable metric type", filter.Type)
	case metric.GaugeType:
		value, err = r.gauge.Get(ctx, filter.Name)
		if err != nil {
			return nil, err
		}
	case metric.CounterType:
		value, err = r.counter.Get(ctx, filter.Name)
		if err != nil {
			return nil, err
		}
	}

	r.Log(ctx).Info().Msg("get metric finished")

	return value, nil
}

func (r repository) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "metric database repository").Logger()

	return &logger
}
