package agent

import (
	"testing"
)

func TestGetAgentMetrics(t *testing.T) {
	t.Run("Should return metrics", func(t *testing.T) {
		if got := GetAgentMetrics(); len(got) != 29 {
			t.Errorf("len GetAgentMetrics() = %v, want %v", len(got), 29)
		}
	})
}
