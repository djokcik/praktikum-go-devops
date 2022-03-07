package metric

import (
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

// Gauge metric type
type Gauge float64

// GaugeType enum value
const GaugeType = "gauge"

// GetLoggerContext zerolog wrapper for metricType
func (o Gauge) GetLoggerContext(metricName string) func(logCtx zerolog.Context) zerolog.Context {
	return func(logCtx zerolog.Context) zerolog.Context {
		return logCtx.
			Str(logging.MetricType, "gauge").
			Str(logging.MetricName, metricName)
	}
}

// GaugeBaseMetric basic counter struct type
type GaugeBaseMetric struct {
}

func (v GaugeBaseMetric) Type() string {
	return GaugeType
}

// GaugeDto dto counter for transmitter and receiver metric
type GaugeDto struct {
	Name  string
	Value Gauge
}
