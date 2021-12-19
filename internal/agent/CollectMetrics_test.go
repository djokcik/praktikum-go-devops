package agent

import (
	"github.com/Jokcik/praktikum-go-devops/internal/agent/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_updateMetrics(f *testing.T) {
	f.Run("Should update mapMetrics from metrics", func(t *testing.T) {
		m := mocks.AgentMetric{Mock: mock.Mock{}}
		m.On("Name").Return("TestName")
		m.On("GetValue").Return("TestValue")

		updatedMetrics := make(map[string]SendAgentMetric)
		metrics := []AgentMetric{&m}

		updateMetrics(updatedMetrics, metrics)

		m.AssertNumberOfCalls(t, "Name", 1)
		m.AssertNumberOfCalls(t, "GetValue", 1)
		require.Equal(t, updatedMetrics["TestName"], SendAgentMetric{Metric: metrics[0], Value: "TestValue"})
	})
}
