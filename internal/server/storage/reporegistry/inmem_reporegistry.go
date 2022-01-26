package reporegistry

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/helpers"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/counterrepo"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/gaugerepo"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/storer"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/storer/filestorer"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
	"sync"
)

type inmemRepoRegistry struct {
	counterRepo counterrepo.Repository
	gaugeRepo   gaugerepo.Repository
}

// NewInMem .
func NewInMem(ctx context.Context, wg *sync.WaitGroup, cfg server.Config) RepoRegistry {
	r := &inmemRepoRegistry{}

	var store storer.MetricStorer

	if cfg.StoreFile != "" {
		store = filestorer.NewMetricFileStorer(ctx, cfg)
	} else {
		r.Log(ctx).Info().Msg("saving metrics to the store is disabled")
	}

	r.counterRepo = counterrepo.NewInMem(store, cfg)
	r.gaugeRepo = gaugerepo.NewInMem(store, cfg)

	if store != nil {
		wg.Add(1)
		go func() {
			<-ctx.Done()
			defer wg.Done()

			store.SaveDBValue(ctx)
			store.Close()
		}()

		if cfg.Restore {
			store.RestoreDBValue(ctx)
		}

		if cfg.StoreInterval != 0 {
			go helpers.SetTicker(func() { store.SaveDBValue(ctx) }, cfg.StoreInterval)
		}
	}

	return r
}

func (r *inmemRepoRegistry) GetCounterRepo() counterrepo.Repository {
	return r.counterRepo
}

func (r *inmemRepoRegistry) GetGaugeRepo() gaugerepo.Repository {
	return r.gaugeRepo
}

func (r *inmemRepoRegistry) Ping(ctx context.Context) error {
	return nil
}

func (r *inmemRepoRegistry) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "inmemRepoRegistry").Logger()

	return &logger
}
