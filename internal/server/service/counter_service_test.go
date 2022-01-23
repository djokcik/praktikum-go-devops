package service

import (
	"context"
	"errors"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCounterServiceImpl_GetOne(t *testing.T) {
	t.Run("1. Should return metric", func(t *testing.T) {
		m := mocks.MetricRepository{Mock: mock.Mock{}}
		m.On("Get", context.Background(), storage.GetRepositoryFilter{
			Name: "TestMetric",
			Type: metric.CounterType,
		}).Return(metric.Counter(123), nil)

		service := CounterServiceImpl{Repo: &m}

		val, err := service.GetOne(context.Background(), "TestMetric")

		require.Equal(t, err, nil)
		require.Equal(t, val, metric.Counter(123))
	})

	t.Run("2. Should return error when repo return error", func(t *testing.T) {
		m := mocks.MetricRepository{Mock: mock.Mock{}}
		m.On("Get", context.Background(), storage.GetRepositoryFilter{
			Name: "TestMetric",
			Type: metric.CounterType,
		}).Return(nil, errors.New("testError"))

		service := CounterServiceImpl{Repo: &m}

		val, err := service.GetOne(context.Background(), "TestMetric")

		require.Equal(t, err, errors.New("testError"))
		require.Equal(t, val, metric.Counter(0))
	})
}

func TestCounterServiceImpl_Update(t *testing.T) {
	t.Run("1. Should update metric", func(t *testing.T) {
		m := mocks.MetricRepository{Mock: mock.Mock{}}
		m.On("Update", context.Background(), "TestMetric", metric.Counter(123)).Return(true, nil)

		service := CounterServiceImpl{Repo: &m}

		val, err := service.Update(context.Background(), "TestMetric", metric.Counter(123))

		m.AssertNumberOfCalls(t, "Update", 1)
		require.Equal(t, err, nil)
		require.Equal(t, val, true)
	})
}

func TestCounterServiceImpl_List(t *testing.T) {
	t.Run("1. Should return list metrics", func(t *testing.T) {
		metricList := []metric.Metric{{Name: "TestType", Value: "TestValue"}}

		m := mocks.MetricRepository{Mock: mock.Mock{}}
		m.On("List", context.Background(), storage.ListRepositoryFilter{Type: metric.CounterType}).Return(metricList, nil)

		service := CounterServiceImpl{Repo: &m}

		metrics, err := service.List(context.Background())

		m.AssertNumberOfCalls(t, "List", 1)
		require.Equal(t, err, nil)
		require.Equal(t, metrics, metricList)
	})
}

func TestCounterServiceImpl_Increase(t *testing.T) {
	t.Run("1. Should add value to old metric", func(t *testing.T) {
		m := mocks.MetricRepository{Mock: mock.Mock{}}
		m.On("Update", context.Background(), "TestMetric", metric.Counter(125)).Return(true, nil)
		m.On("Get", context.Background(), storage.GetRepositoryFilter{
			Name: "TestMetric",
			Type: metric.CounterType,
		}).Return(metric.Counter(100), nil)

		service := CounterServiceImpl{Repo: &m}
		err := service.Increase(context.Background(), "TestMetric", 25)

		m.AssertNumberOfCalls(t, "Update", 1)
		m.AssertNumberOfCalls(t, "Get", 1)
		require.Equal(t, err, nil)
	})

	t.Run("2. Should update metric if didn`t find value in state", func(t *testing.T) {
		m := mocks.MetricRepository{Mock: mock.Mock{}}
		m.On("Update", context.Background(), "TestMetric", metric.Counter(25)).Return(true, nil)
		m.On("Get", context.Background(), storage.GetRepositoryFilter{
			Name: "TestMetric",
			Type: metric.CounterType,
		}).Return(metric.Counter(0), storage.ErrValueNotFound)

		service := CounterServiceImpl{Repo: &m}
		err := service.Increase(context.Background(), "TestMetric", 25)

		m.AssertNumberOfCalls(t, "Update", 1)
		m.AssertNumberOfCalls(t, "Get", 1)
		require.Equal(t, err, nil)
	})

	t.Run("3. Should update metric if find value in state return error", func(t *testing.T) {
		m := mocks.MetricRepository{Mock: mock.Mock{}}
		m.On("Update", context.Background(), mock.Anything, mock.Anything).Return(true, nil)
		m.On("Get", context.Background(), storage.GetRepositoryFilter{
			Name: "TestMetric",
			Type: metric.CounterType,
		}).Return(metric.Counter(0), errors.New("TestError"))

		service := CounterServiceImpl{Repo: &m}
		err := service.Increase(context.Background(), "TestMetric", 25)

		require.Equal(t, err, errors.New("TestError"))
		m.AssertNumberOfCalls(t, "Get", 1)
		m.AssertNumberOfCalls(t, "Update", 0)
	})
}
