package agent

import (
	"testing"
)

func TestGetAgentMetrics(t *testing.T) {
	t.Run("Should return 27 metrics", func(t *testing.T) {
		if got := GetAgentMetrics(); len(got) != 27 {
			t.Errorf("len GetAgentMetrics() = %v, want %v", len(got), 27)
		}
	})
}
