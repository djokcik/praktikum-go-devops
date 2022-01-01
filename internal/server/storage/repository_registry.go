package storage

import (
	"context"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/model"
	"reflect"
	"sync"
)

//go:generate mockery --name=Repository

type GetRepositoryFilter struct {
	Type string
	Name string
}

type ListRepositoryFilter struct {
	Type string
}

type Repository interface {
	Configure(ctx context.Context, wg *sync.WaitGroup, db *model.Database, cfg *server.Config)
	Update(id interface{}, entity interface{}) (bool, error)
	List(filter *ListRepositoryFilter) (interface{}, error)
	Get(filter *GetRepositoryFilter) (interface{}, error)
}

type BaseRepository struct {
	db *model.Database
}

func (r *BaseRepository) Configure(ctx context.Context, wg *sync.WaitGroup, db *model.Database, cfg *server.Config) {
	r.db = db
	db.CounterMapMetric = make(map[string]metric.Counter)
	db.GaugeMapMetric = make(map[string]metric.Gauge)
}

type RepositoryRegistry struct {
	registry map[string]Repository

	db  *model.Database
	cfg *server.Config
}

func (r *RepositoryRegistry) registerRepositories(ctx context.Context, wg *sync.WaitGroup, repositories []Repository) {
	for _, repository := range repositories {
		repositoryName := reflect.TypeOf(repository).Elem().Name()
		repository.Configure(ctx, wg, r.db, r.cfg)
		r.registry[repositoryName] = repository
	}
}

func NewRepositoryRegistry(ctx context.Context, wg *sync.WaitGroup, cfg *server.Config, db *model.Database, repository ...Repository) *RepositoryRegistry {
	r := &RepositoryRegistry{
		registry: map[string]Repository{},
		db:       db,
		cfg:      cfg,
	}

	r.registerRepositories(ctx, wg, repository)
	return r
}

func (r *RepositoryRegistry) Repository(repositoryName string) (Repository, error) {
	if repository, ok := r.registry[repositoryName]; ok {
		return repository, nil
	}

	return nil, fmt.Errorf("repository %s does not exist", repositoryName)
}
