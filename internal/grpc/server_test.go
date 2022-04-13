package grpc

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/service/mocks"
	pb "github.com/djokcik/praktikum-go-devops/pkg/proto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetricsServer_SendMetric(t *testing.T) {
	t.Run("Should update counter and gauge list", func(t *testing.T) {
		mCounter := mocks.CounterService{Mock: mock.Mock{}}
		mCounter.On("Verify", mock.Anything, "TestCounterMetric", metric.Counter(10), "hashCounter").
			Return(true)
		mCounter.On("UpdateList", mock.Anything, []metric.CounterDto{{Name: "TestCounterMetric", Value: metric.Counter(10)}}).
			Return(nil)

		mGauge := mocks.GaugeService{Mock: mock.Mock{}}
		mGauge.On("Verify", mock.Anything, "TestGaugeMetric", metric.Gauge(1.5), "hashGauge").
			Return(true)
		mGauge.On("UpdateList", mock.Anything, []metric.GaugeDto{{Name: "TestGaugeMetric", Value: metric.Gauge(1.5)}}).
			Return(nil)

		// `[{"id":"TestCounterMetric","type":"counter","delta":10,"hash":"hashCounter"},{"id":"TestGaugeMetric","type":"gauge","value":1.5,"hash":"hashGauge"}]`
		server := &MetricsServer{Counter: &mCounter, Gauge: &mGauge}
		res, err := server.SendMetric(context.Background(), &pb.SendMetricRequest{Metrics: []*pb.Metric{
			{ID: "TestCounterMetric", MType: "counter", Delta: 10, Hash: "hashCounter"},
			{ID: "TestGaugeMetric", MType: "gauge", Value: 1.5, Hash: "hashGauge"},
		}})

		require.Equal(t, err, nil)
		require.Equal(t, res, &pb.SendMetricResponse{})

		mCounter.AssertNumberOfCalls(t, "Verify", 1)
		mCounter.AssertNumberOfCalls(t, "UpdateList", 1)
		mGauge.AssertNumberOfCalls(t, "Verify", 1)
		mGauge.AssertNumberOfCalls(t, "UpdateList", 1)
	})
}
