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

//go:generate mockery --name=counterDatabaseService --exported

type counterDatabaseService interface {
	Get(ctx context.Context, name string) (metric.Counter, error)
	List(ctx context.Context) ([]metric.Metric, error)
	Update(ctx context.Context, name string, value metric.Counter) error
}

type counterDatabaseServiceImpl struct {
	db *sql.DB
}

func (c counterDatabaseServiceImpl) Get(ctx context.Context, name string) (metric.Counter, error) {
	row := c.db.QueryRowContext(ctx, "select value from counter_metric where id = $1", name)
	if row.Err() != nil {
		c.Log(ctx).Error().Err(row.Err()).Msg("Get: query return error")
		return 0, row.Err()
	}

	var metricValue metric.Counter
	err := row.Scan(&metricValue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrValueNotFound
		}

		c.Log(ctx).Error().Err(err).Msg("Get: scan return error")
		return 0, err
	}

	return metricValue, nil
}

func (c counterDatabaseServiceImpl) List(ctx context.Context) ([]metric.Metric, error) {
	var metricList []metric.Metric

	rows, err := c.db.QueryContext(ctx, "select id, value from counter_metric")
	if err != nil {
		c.Log(ctx).Error().Err(err).Msg("List: query return error")
		return nil, err
	}

	for rows.Next() {
		var name string
		var value metric.Counter

		err := rows.Scan(&name, &value)
		if err != nil {
			c.Log(ctx).Error().Err(err).Msg("List: scan return error")
			return nil, err
		}

		metricList = append(metricList, metric.Metric{Name: name, Value: value})
	}

	if rows.Err() != nil {
		c.Log(ctx).Error().Err(err).Msg("List: query rows was error")
		return nil, err
	}

	return metricList, nil
}

func (c counterDatabaseServiceImpl) Update(ctx context.Context, name string, value metric.Counter) error {
	query := `INSERT INTO counter_metric(id, value) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET value = excluded.value`
	_, err := c.db.ExecContext(ctx, query, name, value)
	if err != nil {
		c.Log(ctx).Error().Err(err).Msgf("don`t save counter metric %s with value %v", name, value)
		return err
	}

	return nil
}

func (c counterDatabaseServiceImpl) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "counter database repository").Logger()

	return &logger
}
