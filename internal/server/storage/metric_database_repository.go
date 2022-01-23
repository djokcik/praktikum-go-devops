package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	_ "github.com/jackc/pgx/v4/stdlib"
	"os"
)

type MetricDatabaseRepository struct {
	db *sql.DB
}

func NewMetricDatabaseRepository(ctx context.Context, cfg server.Config) MetricRepository {
	r := &MetricDatabaseRepository{}

	db, err := sql.Open("pgx", cfg.DatabaseDsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	r.db = db

	go func() {
		<-ctx.Done()
		r.db.Close()
	}()

	return r
}

func (r MetricDatabaseRepository) Ping() error {
	return r.db.Ping()
}

func (r MetricDatabaseRepository) Update(ctx context.Context, name string, value interface{}) (bool, error) {
	return true, nil
}

func (r MetricDatabaseRepository) List(ctx context.Context, filter ListRepositoryFilter) (interface{}, error) {
	var metricList []metric.Metric

	return metricList, nil
}

func (r MetricDatabaseRepository) Get(ctx context.Context, filter GetRepositoryFilter) (interface{}, error) {
	return nil, nil
}
