package metricdatabase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=gaugeDatabaseService --exported

type gaugeDatabaseService interface {
	Get(ctx context.Context, name string) (metric.Gauge, error)
	List(ctx context.Context) ([]metric.Metric, error)
	Update(ctx context.Context, name string, value metric.Gauge) error
}

type gaugeDatabaseServiceImpl struct {
	db *sql.DB
}

func (g gaugeDatabaseServiceImpl) Get(ctx context.Context, name string) (metric.Gauge, error) {
	row := g.db.QueryRowContext(ctx, "select value from gauge_metric where id = $1", name)
	if row.Err() != nil {
		g.Log(ctx).Error().Err(row.Err()).Msg("Get: query return error")
		return metric.Gauge(0), row.Err()
	}

	var metricValue metric.Gauge
	err := row.Scan(&metricValue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrValueNotFound
		}

		g.Log(ctx).Error().Err(err).Msg("Get: scan return error")
		return metric.Gauge(0), err
	}

	return metricValue, nil
}

func (g gaugeDatabaseServiceImpl) List(ctx context.Context) ([]metric.Metric, error) {
	var metricList []metric.Metric

	rows, err := g.db.QueryContext(ctx, "select id, value from gauge_metric")
	if err != nil {
		g.Log(ctx).Error().Err(err).Msg("List: query return error")
		return nil, err
	}

	for rows.Next() {
		var name string
		var value metric.Gauge

		err := rows.Scan(&name, &value)
		if err != nil {
			g.Log(ctx).Error().Err(err).Msg("List: scan return error")
			return nil, err
		}

		metricList = append(metricList, metric.Metric{Name: name, Value: value})
	}

	if rows.Err() != nil {
		g.Log(ctx).Error().Err(err).Msg("List: query rows was error")
		return nil, err
	}

	return metricList, nil
}

func (g gaugeDatabaseServiceImpl) Update(ctx context.Context, name string, value metric.Gauge) error {
	query := `INSERT INTO gauge_metric(id, value) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET value = excluded.value;`
	_, err := g.db.ExecContext(ctx, query, name, value)
	if err != nil {
		g.Log(ctx).Error().Err(err).Msgf("don`t save gauge metric %s with value %v", name, value)
		return err
	}

	return nil
}

func (g gaugeDatabaseServiceImpl) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "gauge database repository").Logger()

	return &logger
}
