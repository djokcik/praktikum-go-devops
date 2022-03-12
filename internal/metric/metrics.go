package metric

import (
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=Metric

// Metric struct with name and value any of metric
type Metric struct {
	Name  string
	Value interface{}
}

// MetricDto for request agent to server
type MetricDto struct {
	// metric name
	ID string `json:"id"`
	// metric type
	MType string `json:"type"`
	// value of metric if metric type is counter
	Delta *int64 `json:"delta,omitempty"`
	// value of metric if metric type is gauge
	Value *float64 `json:"value,omitempty"`
	// hash from metric type and value
	Hash string `json:"hash,omitempty"`
}

func (o *MetricDto) GetLoggerContext(logCtx zerolog.Context) zerolog.Context {
	return logCtx.
		Str(logging.MetricType, o.MType).
		Str(logging.MetricName, o.ID)
}
