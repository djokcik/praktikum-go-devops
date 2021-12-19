package storage

import (
	"fmt"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage/model"
	"reflect"
)

//go:generate mockery --name=Repository

type Repository interface {
	Configure(db *model.Database)
	Update(id interface{}, entity interface{}) (bool, error)
}

type BaseRepository struct {
	db *model.Database
}

func (r *BaseRepository) Configure(db *model.Database) {
	r.db = db
	db.Table = make(map[string]interface{})
}

type RepositoryRegistry struct {
	registry map[string]Repository

	db *model.Database
}

func (r *RepositoryRegistry) registerRepositories(repositories []Repository) {
	for _, repository := range repositories {
		repositoryName := reflect.TypeOf(repository).Elem().Name()
		repository.Configure(r.db)
		r.registry[repositoryName] = repository
	}
}

func NewRepositoryRegistry(db *model.Database, repository ...Repository) *RepositoryRegistry {
	r := &RepositoryRegistry{
		registry: map[string]Repository{},

		db: db,
	}

	r.registerRepositories(repository)
	return r
}

func (r *RepositoryRegistry) Repository(repositoryName string) (Repository, error) {
	if repository, ok := r.registry[repositoryName]; ok {
		return repository, nil
	}

	return nil, fmt.Errorf("repository %s does not exist", repositoryName)
}
