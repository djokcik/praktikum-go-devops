package counterrepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/storageconst"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

type postgresqlRepository struct {
	db *sql.DB
}

// NewPostgreSQL returns a new instance of Repository
func NewPostgreSQL(db *sql.DB) Repository {
	return &postgresqlRepository{
		db: db,
	}
}

func (c *postgresqlRepository) Get(ctx context.Context, name string) (metric.Counter, error) {
	row := c.db.QueryRowContext(ctx, "select value from counter_metric where id = $1", name)
	if row.Err() != nil {
		c.Log(ctx).Error().Err(row.Err()).Msg("Get: query return error")
		return 0, row.Err()
	}

	var metricValue metric.Counter
	err := row.Scan(&metricValue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storageconst.ErrValueNotFound
		}

		c.Log(ctx).Error().Err(err).Msg("Get: scan return error")
		return 0, err
	}

	return metricValue, nil
}

func (c *postgresqlRepository) List(ctx context.Context) ([]metric.Metric, error) {
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

func (c *postgresqlRepository) UpdateList(ctx context.Context, metrics []metric.CounterDto) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		c.Log(ctx).Error().Err(err).Msgf("erro start transaction")
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO counter_metric(id, value) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET value = excluded.value")
	if err != nil {
		c.Log(ctx).Error().Err(err).Msgf("error Prepare transaction")
		return err
	}

	for _, dto := range metrics {
		if _, err := stmt.ExecContext(ctx, dto.Name, dto.Value); err != nil {
			c.Log(ctx).Error().Err(err).Msgf("don`t save counter metric %s with value %v", dto.Name, dto.Value)
			if err = tx.Rollback(); err != nil {
				c.Log(ctx).Error().Err(err).Msgf("update drivers: unable to rollback")
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		c.Log(ctx).Error().Err(err).Msgf("update drivers: unable to commit")
		return err
	}

	return nil
}

func (c *postgresqlRepository) Update(ctx context.Context, name string, value metric.Counter) error {
	query := `INSERT INTO counter_metric(id, value) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET value = excluded.value`
	_, err := c.db.ExecContext(ctx, query, name, value)
	if err != nil {
		c.Log(ctx).Error().Err(err).Msgf("don`t save counter metric %s with value %v", name, value)
		return err
	}

	return nil
}

func (c *postgresqlRepository) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "counter database repository").Logger()

	return &logger
}
