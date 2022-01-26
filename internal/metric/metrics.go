package metric

import (
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=Metric

type Metric struct {
	Name  string
	Value interface{}
}

type MetricDto struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func (o *MetricDto) GetLoggerContext(logCtx zerolog.Context) zerolog.Context {
	return logCtx.
		Str(logging.MetricType, o.MType).
		Str(logging.MetricName, o.ID)
}
