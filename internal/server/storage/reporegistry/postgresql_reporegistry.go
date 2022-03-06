package reporegistry

import (
	"context"
	"database/sql"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/counterrepo"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/gaugerepo"
)

var (
	_ RepoRegistry = (*postgresqlRepoRegistry)(nil)
)

type postgresqlRepoRegistry struct {
	db *sql.DB
}

// NewPostgreSQL .
func NewPostgreSQL(_ context.Context, db *sql.DB) RepoRegistry {
	return &postgresqlRepoRegistry{
		db: db,
	}
}

func (r postgresqlRepoRegistry) GetCounterRepo() counterrepo.Repository {
	return counterrepo.NewPostgreSQL(r.db)
}

func (r postgresqlRepoRegistry) GetGaugeRepo() gaugerepo.Repository {
	return gaugerepo.NewPostgreSQL(r.db)
}

func (r postgresqlRepoRegistry) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
