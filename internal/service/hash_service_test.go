package service

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHashServiceImpl_GetCounterHash(t *testing.T) {
	t.Run("1. Should return hash for Counter", func(t *testing.T) {
		service := NewHashService("test")

		counterHash := service.GetCounterHash(context.Background(), "TestMetric", metric.Counter(10))
		actualHash := service.GetHash(context.Background(), "TestMetric:counter:10")

		require.Equal(t, counterHash, actualHash)
	})

	t.Run("2. Should return hash for Gauge", func(t *testing.T) {
		service := NewHashService("test")

		gaugeHash := service.GetGaugeHash(context.Background(), "TestMetric", metric.Gauge(0.123))
		actualHash := service.GetHash(context.Background(), "TestMetric:gauge:0.123000")

		require.Equal(t, gaugeHash, actualHash)
	})

	t.Run("3. Should verify hash is truth", func(t *testing.T) {
		service := NewHashService("test")

		hash := service.GetHash(context.Background(), "123")
		hash2 := service.GetHash(context.Background(), "123")

		result := service.Verify(context.Background(), hash, hash2)

		require.Equal(t, result, true)
	})

	t.Run("4. Should verify hash is falsy", func(t *testing.T) {
		service := NewHashService("test")

		hash := service.GetHash(context.Background(), "123")
		hash2 := service.GetHash(context.Background(), "124")

		result := service.Verify(context.Background(), hash, hash2)

		require.Equal(t, result, false)
	})
}
