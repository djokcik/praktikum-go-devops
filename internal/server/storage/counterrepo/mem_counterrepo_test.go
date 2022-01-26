package counterrepo

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/storageconst"
	"github.com/stretchr/testify/require"
	"testing"
)

func newTestInMem(ctx context.Context, db map[string]metric.Counter, cfg server.Config) *inmemRepository {
	return &inmemRepository{
		db:  inmemDB{data: db},
		cfg: cfg,
	}
}

func Test_inmemRepository_UpdateList(t *testing.T) {
	t.Run("1. Should return list Counter metrics", func(t *testing.T) {
		db := make(map[string]metric.Counter)

		repository := newTestInMem(context.Background(), db, server.Config{})

		db["TestCounterName"] = metric.Counter(123)

		list, err := repository.List(context.Background())
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, list, []metric.Metric{{Name: "TestCounterName", Value: metric.Counter(123)}})
	})
}

func Test_inmemRepository_Update(t *testing.T) {
	t.Run("1. Should update repository counter map", func(t *testing.T) {
		db := make(map[string]metric.Counter)

		repository := newTestInMem(context.Background(), db, server.Config{})

		err := repository.Update(context.Background(), "MetricName", metric.Counter(123))
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, db["MetricName"], metric.Counter(123))
	})
}

func Test_inmemRepository_List(t *testing.T) {
	t.Run("1. Should return list Counter metrics", func(t *testing.T) {
		db := make(map[string]metric.Counter)

		repository := newTestInMem(context.Background(), db, server.Config{})

		db["TestCounterName"] = metric.Counter(123)

		list, err := repository.List(context.Background())
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, list, []metric.Metric{{Name: "TestCounterName", Value: metric.Counter(123)}})
	})
}

func Test_inmemRepository_Get(t *testing.T) {
	t.Run("1. Should return counter metric", func(t *testing.T) {
		db := make(map[string]metric.Counter)

		repository := newTestInMem(context.Background(), db, server.Config{})

		db["TestCounterName"] = metric.Counter(123)

		val, err := repository.Get(context.Background(), "TestCounterName")
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, val, metric.Counter(123))
	})

	t.Run("2. Should return error when counter metric not found", func(t *testing.T) {
		db := make(map[string]metric.Counter)

		repository := newTestInMem(context.Background(), db, server.Config{})

		val, err := repository.Get(context.Background(), "TestCounterName")

		require.Equal(t, val, metric.Counter(0))
		require.Equal(t, err, storageconst.ErrValueNotFound)
	})
}
