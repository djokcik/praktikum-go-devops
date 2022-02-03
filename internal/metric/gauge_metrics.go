package metric

import (
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

type Gauge float64

const GaugeType = "gauge"

func (o Gauge) GetLoggerContext(metricName string) func(logCtx zerolog.Context) zerolog.Context {
	return func(logCtx zerolog.Context) zerolog.Context {
		return logCtx.
			Str(logging.MetricType, "gauge").
			Str(logging.MetricName, metricName)
	}
}

type GaugeBaseMetric struct {
}

func (v GaugeBaseMetric) Type() string {
	return GaugeType
}

type GaugeDto struct {
	Name  string
	Value Gauge
}
