package storage

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetricRepository_Update(t *testing.T) {
	t.Run("1. Should update repository counter map", func(t *testing.T) {
		db := new(model.Database)

		repository := new(MetricRepository)
		repository.Configure(db)

		val, err := repository.Update("MetricName", metric.Counter(123))
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, val, true)
		require.Equal(t, db.CounterMapMetric["MetricName"], metric.Counter(123))
	})

	t.Run("2. Should update repository gauge map", func(t *testing.T) {
		db := new(model.Database)

		repository := new(MetricRepository)
		repository.Configure(db)

		val, err := repository.Update("MetricName", metric.Gauge(0.123))
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, val, true)
		require.Equal(t, db.GaugeMapMetric["MetricName"], metric.Gauge(0.123))
	})

	t.Run("3. Should return error when metric isn`t available", func(t *testing.T) {
		db := new(model.Database)

		repository := new(MetricRepository)
		repository.Configure(db)

		val, err := repository.Update("MetricName", 123)

		require.Equal(t, val, false)
		require.Equal(t, err == nil, false)
	})
}

func TestMetricRepository_List(t *testing.T) {
	t.Run("1. Should return list Counter metrics", func(t *testing.T) {
		db := new(model.Database)

		repository := new(MetricRepository)
		repository.Configure(db)

		db.GaugeMapMetric["TestGaugeName"] = metric.Gauge(0.123)
		db.CounterMapMetric["TestCounterName"] = metric.Counter(123)

		list, err := repository.List(&ListRepositoryFilter{Type: metric.CounterType})
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, list, []metric.Metric{{Name: "TestCounterName", Value: metric.Counter(123)}})
	})

	t.Run("2. Should return list Gauge metrics", func(t *testing.T) {
		db := new(model.Database)

		repository := new(MetricRepository)
		repository.Configure(db)

		db.GaugeMapMetric["TestGaugeName"] = metric.Gauge(0.123)
		db.CounterMapMetric["TestCounterName"] = metric.Counter(123)

		list, err := repository.List(&ListRepositoryFilter{Type: metric.GaugeType})
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, list, []metric.Metric{{Name: "TestGaugeName", Value: metric.Gauge(0.123)}})
	})
}

func TestMetricRepository_Get(t *testing.T) {
	t.Run("1. Should return counter metric", func(t *testing.T) {
		db := new(model.Database)

		repository := new(MetricRepository)
		repository.Configure(db)

		db.GaugeMapMetric["TestGaugeName"] = metric.Gauge(0.123)
		db.CounterMapMetric["TestCounterName"] = metric.Counter(123)

		val, err := repository.Get(&GetRepositoryFilter{Type: metric.CounterType, Name: "TestCounterName"})
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, val, metric.Counter(123))
	})

	t.Run("2. Should return error when counter metric not found", func(t *testing.T) {
		db := new(model.Database)

		repository := new(MetricRepository)
		repository.Configure(db)

		val, err := repository.Get(&GetRepositoryFilter{
			Type: metric.CounterType,
			Name: "TestCounterName",
		})

		require.Equal(t, val, 0)
		require.Equal(t, err.Error(), ValueNotFound)
	})

	t.Run("2. Should return gauge metric", func(t *testing.T) {
		db := new(model.Database)

		repository := new(MetricRepository)
		repository.Configure(db)

		db.GaugeMapMetric["TestGaugeName"] = metric.Gauge(0.123)
		db.CounterMapMetric["TestCounterName"] = metric.Counter(123)

		val, err := repository.Get(&GetRepositoryFilter{Type: metric.GaugeType, Name: "TestGaugeName"})
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, val, metric.Gauge(0.123))
	})

	t.Run("3. Should return default gauge metric", func(t *testing.T) {
		db := new(model.Database)

		repository := new(MetricRepository)
		repository.Configure(db)

		val, err := repository.Get(&GetRepositoryFilter{
			Type: metric.GaugeType,
			Name: "TestGaugeName",
		})

		require.Equal(t, val, 0)
		require.Equal(t, err.Error(), ValueNotFound)
	})
}
