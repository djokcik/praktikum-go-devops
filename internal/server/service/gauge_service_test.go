package service

import (
	"errors"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGaugeServiceImpl_Update(t *testing.T) {
	t.Run("1. Should update metric", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", "TestMetric", metric.Gauge(0.123)).Return(true, nil)

		service := GaugeServiceImpl{Repo: &m}

		val, err := service.Update("TestMetric", metric.Gauge(0.123))

		m.AssertNumberOfCalls(t, "Update", 1)
		require.Equal(t, err, nil)
		require.Equal(t, val, true)
	})
}

func TestGaugeServiceImpl_GetOne(t *testing.T) {
	t.Run("1. Should return metric", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Get", &storage.GetRepositoryFilter{
			Name: "TestMetric",
			Type: metric.GaugeType,
		}).Return(metric.Gauge(0.123), nil)

		service := GaugeServiceImpl{Repo: &m}

		val, err := service.GetOne("TestMetric")

		require.Equal(t, err, nil)
		require.Equal(t, val, metric.Gauge(0.123))
	})

	t.Run("2. Should return error when repo return error", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Get", &storage.GetRepositoryFilter{
			Name: "TestMetric",
			Type: metric.GaugeType,
		}).Return(nil, errors.New("testError"))

		service := GaugeServiceImpl{Repo: &m}

		val, err := service.GetOne("TestMetric")

		require.Equal(t, err, errors.New("testError"))
		require.Equal(t, val, metric.Gauge(0))
	})
}

func TestGaugeServiceImpl_List(t *testing.T) {
	t.Run("1. Should return list metrics", func(t *testing.T) {
		metricList := []metric.Metric{{Name: "TestType", Value: "TestValue"}}

		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("List", &storage.ListRepositoryFilter{Type: metric.GaugeType}).Return(metricList, nil)

		service := GaugeServiceImpl{Repo: &m}

		metrics, err := service.List()

		m.AssertNumberOfCalls(t, "List", 1)
		require.Equal(t, err, nil)
		require.Equal(t, metrics, metricList)
	})
}