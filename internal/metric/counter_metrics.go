// Package metric is general information about metrics
package metric

import (
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

// Counter metric type
type Counter int64

// CounterType enum value
const CounterType = "counter"

// GetLoggerContext zerolog wrapper for metricType
func (c Counter) GetLoggerContext(metricName string) func(logCtx zerolog.Context) zerolog.Context {
	return func(logCtx zerolog.Context) zerolog.Context {
		return logCtx.
			Str(logging.MetricType, "counter").
			Str(logging.MetricName, metricName)
	}
}

// CounterBaseMetric basic counter struct type
type CounterBaseMetric struct {
}

func (v CounterBaseMetric) Type() string {
	return CounterType
}

// CounterDto dto counter for transmitter and receiver metric
type CounterDto struct {
	Name  string
	Value Counter
}
