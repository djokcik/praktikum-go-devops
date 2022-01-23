package store

import (
	"context"
	"database/sql"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/model"
	"github.com/djokcik/praktikum-go-devops/migration"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
)

type MetricDbStorer struct {
	inMemoryDB *model.InMemoryMetricDB
	cfg        server.Config

	db *sql.DB
}

func NewMetricDbStorer(ctx context.Context, inMemoryDB *model.InMemoryMetricDB, cfg server.Config) MetricStorer {
	s := &MetricDbStorer{inMemoryDB: inMemoryDB, cfg: cfg}

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

func (s *MetricDbStorer) RestoreDBValue(ctx context.Context) {

}

func (s MetricDbStorer) SaveDBValue(ctx context.Context) {

}

func (s MetricDbStorer) Ping(ctx context.Context) error {
	return s.db.Ping()
}

func (s *MetricDbStorer) Close() {
	s.db.Close()
}

func (s *MetricDbStorer) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "metric db storer").Logger()

	return &logger
}
