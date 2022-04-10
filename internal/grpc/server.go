package grpc

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/service"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	pb "github.com/djokcik/praktikum-go-devops/pkg/proto"
	"github.com/rs/zerolog"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer
	Counter service.CounterService
	Gauge   service.GaugeService
}

func (s MetricsServer) SendMetric(ctx context.Context, in *pb.SendMetricRequest) (*pb.SendMetricResponse, error) {
	logger := s.Log(ctx).With().Str(logging.ServiceKey, "gRPC SendMetric").Logger()
	ctx = logging.SetCtxLogger(ctx, logger)

	metricsDto := in.Metrics

	counterMetrics := make([]metric.CounterDto, 0)
	for _, metricDto := range metricsDto {
		if metricDto.MType == metric.CounterType {
			name := metricDto.ID
			value := metric.Counter(metricDto.Delta)

			if !s.Counter.Verify(ctx, name, value, metricDto.Hash) {
				continue
			}

			counterMetrics = append(counterMetrics, metric.CounterDto{Name: name, Value: value})
		}
	}

	gaugeMetrics := make([]metric.GaugeDto, 0)
	for _, metricDto := range metricsDto {
		if metricDto.MType == metric.GaugeType {
			name := metricDto.ID
			value := metric.Gauge(metricDto.Value)

			if !s.Gauge.Verify(ctx, name, value, metricDto.Hash) {
				continue
			}

			gaugeMetrics = append(gaugeMetrics, metric.GaugeDto{Name: name, Value: value})
		}
	}

	if len(counterMetrics) != 0 {
		err := s.Counter.UpdateList(ctx, counterMetrics)
		if err != nil {
			s.Log(ctx).Error().Err(err).Msg("invalid save counter metrics")
			return &pb.SendMetricResponse{Error: "invalid save counter metrics"}, nil
		}
	}

	if len(gaugeMetrics) != 0 {
		err := s.Gauge.UpdateList(ctx, gaugeMetrics)
		if err != nil {
			s.Log(ctx).Error().Err(err).Msg("invalid save gauge metrics")
			return &pb.SendMetricResponse{Error: "invalid save gauge metrics"}, nil
		}
	}

	s.Log(ctx).Info().Msg("json update list metrics handled")

	return &pb.SendMetricResponse{}, nil
}

func (s *MetricsServer) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "MetricsServer").Logger()

	return &logger
}
