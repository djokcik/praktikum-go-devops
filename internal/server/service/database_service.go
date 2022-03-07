package service

import (
	"context"
	"database/sql"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/migration"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"sync"
)

//go:generate mockery --name=DatabaseService

// DatabaseService - iterface for database connection
type DatabaseService interface {
	Open(ctx context.Context, wg *sync.WaitGroup) (*sql.DB, error)
}

type databaseService struct {
	cfg server.Config
}

func NewDatabaseService(_ context.Context, cfg server.Config) DatabaseService {
	return &databaseService{cfg: cfg}
}

// Open - open database with migration
func (r databaseService) Open(ctx context.Context, wg *sync.WaitGroup) (*sql.DB, error) {
	db, err := sql.Open("pgx", r.cfg.DatabaseDsn)
	if err != nil {
		r.Log(ctx).Fatal().Err(err).Msgf("Unable to connect to database")
		return nil, err
	}

	err = migration.CreateCounterTable(db)
	if err != nil {
		r.Log(ctx).Warn().Err(err).Msgf("couldn't create counter table")
	}

	err = migration.CreateGaugeTable(db)
	if err != nil {
		r.Log(ctx).Warn().Err(err).Msgf("couldn't create gauge table")
	}

	wg.Add(1)
	go func() {
		<-ctx.Done()
		defer wg.Done()

		db.Close()
	}()

	return db, nil
}

func (r databaseService) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "metric database repository").Logger()

	return &logger
}
