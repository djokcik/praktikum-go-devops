package servermetrics

import (
	"testing"
)

func TestGetServerGaugeMetrics(t *testing.T) {
	t.Run("Should return 26 gauge metrics", func(t *testing.T) {
		if got := GetServerGaugeMetrics(); len(got) != 26 {
			t.Errorf("len GetServerGaugeMetrics() = %v, want %v", len(got), 26)
		}
	})
}

func TestGetServerCounterMetrics(t *testing.T) {
	t.Run("Should return 1 counter metrics", func(t *testing.T) {
		if got := GetServerCounterMetrics(); len(got) != 1 {
			t.Errorf("len GetServerGaugeMetrics() = %v, want %v", len(got), 1)
		}
	})
}
