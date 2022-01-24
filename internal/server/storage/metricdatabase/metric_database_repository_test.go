package metricdatabase

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/metricdatabase/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_repository_Ping(t *testing.T) {
	t.Run("1. should return error with ping", func(t *testing.T) {
		db, dbMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		r := &repository{db: db}

		expectErr := errors.New("TestError")
		dbMock.ExpectPing().WillReturnError(expectErr)

		err = r.Ping(context.Background())

		require.Equal(t, err, expectErr)
	})
}

func Test_repository_Update(t *testing.T) {
	t.Run("1. Should update repository counter map", func(t *testing.T) {
		m := &mocks.CounterDatabaseService{Mock: mock.Mock{}}

		repository := &repository{counter: m}
		m.On("Update", mock.Anything, "MetricName", metric.Counter(123)).Return(nil)

		val, err := repository.Update(context.Background(), "MetricName", metric.Counter(123))

		require.Equal(t, val, true)
		require.Equal(t, err, nil)
	})

	t.Run("2. Should update repository gauge map", func(t *testing.T) {
		m := &mocks.GaugeDatabaseService{Mock: mock.Mock{}}

		repository := &repository{gauge: m}
		m.On("Update", mock.Anything, "MetricName", metric.Gauge(0.123)).Return(nil)

		val, err := repository.Update(context.Background(), "MetricName", metric.Gauge(0.123))
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, val, true)
		require.Equal(t, err, nil)
	})

	t.Run("3. Should return error when metric isn`t available", func(t *testing.T) {
		repository := &repository{}

		val, err := repository.Update(context.Background(), "MetricName", 123)

		require.Equal(t, val, false)
		require.NotEqual(t, err, nil)
	})
}

func Test_repository_List(t *testing.T) {
	t.Run("1. Should return list Counter metrics", func(t *testing.T) {
		m := &mocks.CounterDatabaseService{Mock: mock.Mock{}}

		repository := &repository{counter: m}
		m.On("List", mock.Anything).Return([]metric.Metric{{Name: "TestCounterName", Value: metric.Counter(123)}}, nil)

		list, err := repository.List(context.Background(), storage.ListRepositoryFilter{Type: metric.CounterType})
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, list, []metric.Metric{{Name: "TestCounterName", Value: metric.Counter(123)}})
	})

	t.Run("2. Should return list Gauge metrics", func(t *testing.T) {
		m := &mocks.GaugeDatabaseService{Mock: mock.Mock{}}

		repository := &repository{gauge: m}
		m.On("List", mock.Anything).Return([]metric.Metric{{Name: "TestGaugeName", Value: metric.Gauge(0.123)}}, nil)

		list, err := repository.List(context.Background(), storage.ListRepositoryFilter{Type: metric.GaugeType})
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, list, []metric.Metric{{Name: "TestGaugeName", Value: metric.Gauge(0.123)}})
	})
}

func Test_repository_Get(t *testing.T) {
	t.Run("1. Should return counter metric", func(t *testing.T) {
		m := &mocks.CounterDatabaseService{Mock: mock.Mock{}}

		repository := &repository{counter: m}
		m.On("Get", mock.Anything, "TestCounterName").Return(metric.Counter(123), nil)

		val, err := repository.Get(context.Background(), storage.GetRepositoryFilter{Type: metric.CounterType, Name: "TestCounterName"})
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, val, metric.Counter(123))
	})

	t.Run("2. Should return gauge metric", func(t *testing.T) {
		m := &mocks.GaugeDatabaseService{Mock: mock.Mock{}}

		repository := &repository{gauge: m}
		m.On("Get", mock.Anything, "TestGaugeName").Return(metric.Gauge(0.123), nil)

		val, err := repository.Get(context.Background(), storage.GetRepositoryFilter{Type: metric.GaugeType, Name: "TestGaugeName"})
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, val, metric.Gauge(0.123))
	})
}
