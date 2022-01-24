package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/model"
	"github.com/djokcik/praktikum-go-devops/migration"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"strconv"
	"strings"
)

type metricDBStorer struct {
	inMemoryDB *model.InMemoryMetricDB
	cfg        server.Config

	db *sql.DB
}

func NewMetricDBStorer(ctx context.Context, inMemoryDB *model.InMemoryMetricDB, cfg server.Config) MetricStorer {
	s := &metricDBStorer{inMemoryDB: inMemoryDB, cfg: cfg}

	db, err := sql.Open("pgx", cfg.DatabaseDsn)
	if err != nil {
		s.Log(ctx).Fatal().Err(err).Msgf("Unable to connect to database")
		return nil
	}

	s.db = db

	err = migration.CreateCounterTable(db)
	if err != nil {
		s.Log(ctx).Warn().Err(err).Msgf("couldn't create counter table")
	}

	err = migration.CreateGaugeTable(db)
	if err != nil {
		s.Log(ctx).Warn().Err(err).Msgf("couldn't create gauge table")
	}

	return s
}

func (s *metricDBStorer) RestoreDBValue(ctx context.Context) {
	s.Log(ctx).Info().Msg("start restore metrics from db")

	// COUNTER
	{
		rows, err := s.db.Query("SELECT type, value FROM counter_metric")
		if err != nil {
			fmt.Println(err)
		}

		defer rows.Close()

		for rows.Next() {
			var metricName string
			var metricValue int

			err = rows.Scan(&metricName, &metricValue)
			if err != nil {
				fmt.Println(err)
			}

			s.inMemoryDB.CounterMapMetric[metricName] = metric.Counter(metricValue)
		}

		// проверяем на ошибки
		err = rows.Err()
		if err != nil {
			fmt.Println(err)
		}
	}

	// GAUGE
	{
		rows, err := s.db.Query("SELECT type, value FROM gauge_metric")
		if err != nil {
			fmt.Println(err)
		}

		defer rows.Close()

		for rows.Next() {
			var metricName string
			var metricValue float64

			err = rows.Scan(&metricName, &metricValue)
			if err != nil {
				fmt.Println(err)
			}

			s.inMemoryDB.GaugeMapMetric[metricName] = metric.Gauge(metricValue)
		}

		// проверяем на ошибки
		err = rows.Err()
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(s.inMemoryDB.GaugeMapMetric)
}

func (s metricDBStorer) SaveDBValue(ctx context.Context) {
	s.Log(ctx).Info().Msg("start save metrics to db")

	// COUNTER
	if len(s.inMemoryDB.CounterMapMetric) > 0 {
		s.db.Exec("DELETE FROM counter_metric")

		sqlStr := "INSERT INTO counter_metric(id, type, value) VALUES "
		var vals []interface{}
		var inserts []string

		i := 1
		for key, val := range s.inMemoryDB.CounterMapMetric {
			inserts = append(inserts, "(?, ?, ?)")
			vals = append(vals, i, key, int(val))
		}

		sqlStr = sqlStr + strings.Join(inserts, ",")
		sqlStr = ReplaceSQL(sqlStr, "?")

		exec, err := s.db.ExecContext(context.Background(), sqlStr, vals...)
		fmt.Println(exec, err)
	}
	//

	// GAUGE
	if len(s.inMemoryDB.GaugeMapMetric) > 0 {
		s.db.Exec("DELETE FROM gauge_metric")

		sqlStr := "INSERT INTO gauge_metric(id, type, value) VALUES "
		var vals []interface{}
		var inserts []string

		i := 1
		for key, val := range s.inMemoryDB.GaugeMapMetric {
			inserts = append(inserts, "(?, ?, ?)")
			vals = append(vals, i, key, float64(val))
			i += 1
		}

		sqlStr = sqlStr + strings.Join(inserts, ",")
		sqlStr = ReplaceSQL(sqlStr, "?")

		exec, err := s.db.ExecContext(context.Background(), sqlStr, vals...)
		fmt.Println(exec, err)
	}
	//

	s.Log(ctx).Info().Msg("finished save metrics to db")
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func (s metricDBStorer) Ping(ctx context.Context) error {
	return s.db.Ping()
}

func (s *metricDBStorer) Close() {
	s.db.Close()
}

func (s *metricDBStorer) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "metric db storer").Logger()

	return &logger
}
