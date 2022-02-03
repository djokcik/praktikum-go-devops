package agent

import (
	"runtime"
	"testing"
)

func TestGetAgentMetrics(t *testing.T) {
	t.Run("Should return metrics", func(t *testing.T) {
		if got := GetAgentMetrics(); len(got) != 29 {
			t.Errorf("len GetAgentMetrics() = %v, want %v", len(got), 29)
		}
	})
}

func Test_agent_GetAgentPsutilMetrics(t *testing.T) {
	t.Run("Should return psutil metrics", func(t *testing.T) {
		cpus := runtime.NumCPU()

		if got := GetAgentPsutilMetrics(); len(got) != cpus+2 {
			t.Errorf("len GetAgentMetrics() = %v, want %v", len(got), cpus+2)
		}
	})
}
