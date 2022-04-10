package agent

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_agent_CollectPsutilMetrics(t *testing.T) {
	t.Run("Should update mapMetrics from metrics", func(t *testing.T) {
		metricAgent := NewAgent(Config{}, nil).(*agent)

		metricAgent.CollectPsutilMetrics(context.Background())()

		require.Equal(t, len(metricAgent.CollectedMetric) > 0, true)
	})
}
