package reporegistry

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/counterrepo"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/gaugerepo"
)

//go:generate mockery --name=RepoRegistry

// RepoRegistry .
type RepoRegistry interface {
	GetCounterRepo() counterrepo.Repository
	GetGaugeRepo() gaugerepo.Repository
	Ping(ctx context.Context) error
}
