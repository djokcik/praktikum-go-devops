package gaugerepo

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/storageconst"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_postgresqlRepository_Get(t *testing.T) {
	t.Run("1. should return gauge metric ", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		service := &postgresqlRepository{db: db}

		rows := sqlmock.NewRows([]string{"value"}).AddRow(metric.Gauge(0.123))
		mock.ExpectQuery("select value from gauge_metric where id = \\$1").WithArgs("TestGauge").WillReturnRows(rows)

		result, err := service.Get(context.Background(), "TestGauge")

		require.Equal(t, err, nil)
		require.Equal(t, result, metric.Gauge(0.123))
	})

	t.Run("2. should return NotFoundMetric", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		service := &postgresqlRepository{db: db}

		rows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("select value from gauge_metric where id = \\$1").WithArgs("TestGauge").WillReturnRows(rows)

		result, err := service.Get(context.Background(), "TestGauge")

		require.Equal(t, err, storageconst.ErrValueNotFound)
		require.Equal(t, result, metric.Gauge(0))
	})

	t.Run("3. should return NotFoundMetric", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		service := &postgresqlRepository{db: db}

		rows := sqlmock.NewRows([]string{})
		mock.
			ExpectQuery("select value from gauge_metric where id = \\$1").
			WithArgs("TestGauge").
			WillReturnRows(rows)

		result, err := service.Get(context.Background(), "TestGauge")

		require.Equal(t, err, storageconst.ErrValueNotFound)
		require.Equal(t, result, metric.Gauge(0))
	})
}

func Test_postgresqlRepository_List(t *testing.T) {
	t.Run("1. should return list metrics", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		service := &postgresqlRepository{db: db}

		rows := sqlmock.
			NewRows([]string{"id", "value"}).
			AddRow("Test1Gauge", 0.1).
			AddRow("Test2Gauge", 0.2)

		mock.ExpectQuery("select id, value from gauge_metric").WillReturnRows(rows)

		result, err := service.List(context.Background())

		require.Equal(t, err, nil)
		require.Equal(t, result, []metric.Metric{
			{Name: "Test1Gauge", Value: metric.Gauge(0.1)},
			{Name: "Test2Gauge", Value: metric.Gauge(0.2)},
		})
	})
}

func Test_postgresqlRepository_Update(t *testing.T) {
	t.Run("1. should update metric", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		service := &postgresqlRepository{db: db}

		mock.
			ExpectExec("INSERT INTO gauge_metric\\(id, value\\) VALUES \\(\\$1, \\$2\\) ON CONFLICT \\(id\\) DO UPDATE SET value = excluded.value").
			WithArgs("TestGauge", metric.Gauge(0.123)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = service.Update(context.Background(), "TestGauge", metric.Gauge(0.123))

		require.Equal(t, err, nil)
	})
}
