package metricdatabase

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_counterDatabaseService_Get(t *testing.T) {
	t.Run("1. should return counter metric ", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		service := &counterDatabaseServiceImpl{db: db}

		rows := sqlmock.NewRows([]string{"value"}).AddRow(metric.Counter(123))
		mock.ExpectQuery("select value from counter_metric where id = \\$1").WithArgs("TestCounter").WillReturnRows(rows)

		result, err := service.Get(context.Background(), "TestCounter")

		require.Equal(t, err, nil)
		require.Equal(t, result, metric.Counter(123))
	})

	t.Run("2. should return NotFoundMetric", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		service := &counterDatabaseServiceImpl{db: db}

		rows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("select value from counter_metric where id = \\$1").WithArgs("TestCounter").WillReturnRows(rows)

		result, err := service.Get(context.Background(), "TestCounter")

		require.Equal(t, err, storage.ErrValueNotFound)
		require.Equal(t, result, metric.Counter(0))
	})

	t.Run("3. should return NotFoundMetric", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		service := &counterDatabaseServiceImpl{db: db}

		rows := sqlmock.NewRows([]string{})
		mock.
			ExpectQuery("select value from counter_metric where id = \\$1").
			WithArgs("TestCounter").
			WillReturnRows(rows)

		result, err := service.Get(context.Background(), "TestCounter")

		require.Equal(t, err, storage.ErrValueNotFound)
		require.Equal(t, result, metric.Counter(0))
	})
}

func Test_counterDatabaseService_List(t *testing.T) {
	t.Run("1. should return list metrics", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		service := &counterDatabaseServiceImpl{db: db}

		rows := sqlmock.
			NewRows([]string{"id", "value"}).
			AddRow("Test1Counter", 1).
			AddRow("Test2Counter", 2)

		mock.ExpectQuery("select id, value from counter_metric").WillReturnRows(rows)

		result, err := service.List(context.Background())

		require.Equal(t, err, nil)
		require.Equal(t, result, []metric.Metric{
			{Name: "Test1Counter", Value: metric.Counter(1)},
			{Name: "Test2Counter", Value: metric.Counter(2)},
		})
	})
}

func Test_counterDatabaseService_Update(t *testing.T) {
	t.Run("1. should update metric", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		service := &counterDatabaseServiceImpl{db: db}

		mock.
			ExpectExec("INSERT INTO counter_metric\\(id, value\\) VALUES \\(\\$1, \\$2\\) ON CONFLICT \\(id\\) DO UPDATE SET value = excluded.value").
			WithArgs("TestCounter", metric.Counter(123)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = service.Update(context.Background(), "TestCounter", metric.Counter(123))

		require.Equal(t, err, nil)
	})
}
