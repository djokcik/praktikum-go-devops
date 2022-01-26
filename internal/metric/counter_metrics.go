package metric

import (
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

type Counter int64

const CounterType = "counter"

func (c Counter) GetLoggerContext(metricName string) func(logCtx zerolog.Context) zerolog.Context {
	return func(logCtx zerolog.Context) zerolog.Context {
		return logCtx.
			Str(logging.MetricType, "counter").
			Str(logging.MetricName, metricName)
	}
}

type CounterBaseMetric struct {
}

func (v CounterBaseMetric) Type() string {
	return CounterType
}

type CounterDto struct {
	Name  string
	Value Counter
}
