package agent

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/agent/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_agent_CollectPsutilMetrics(t *testing.T) {
	t.Run("Should update mapMetrics from metrics", func(t *testing.T) {
		metricAgent := NewAgent(Config{})

		m := mocks.AgentMetric{Mock: mock.Mock{}}
		m.On("Name").Return("TestName")
		m.On("Type").Return("TestType")
		m.On("GetValue").Return("TestValue")

		collectedMetric := make(map[string]SendAgentMetric)
		metricAgent.metrics = []AgentMetric{&m}
		metricAgent.CollectedMetric = collectedMetric

		metricAgent.CollectPsutilMetrics(context.Background())()

		m.AssertNumberOfCalls(t, "Name", 1)
		m.AssertNumberOfCalls(t, "Type", 1)
		m.AssertNumberOfCalls(t, "GetValue", 1)
		require.Equal(t, collectedMetric["TestName"], SendAgentMetric{Name: "TestName", Type: "TestType", Value: "TestValue"})
	})
}
