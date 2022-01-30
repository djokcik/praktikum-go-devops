package handler

import (
	"encoding/json"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"net/http"
)

func (h *Handler) UpdateListJSONHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "UpdateListJSONHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		var metricsDto []metric.MetricDto
		err := json.NewDecoder(r.Body).Decode(&metricsDto)
		if err != nil {
			h.Log(ctx).Error().Err(err).Msg("invalid save counter metrics")
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		counterMetrics := make([]metric.CounterDto, 0)
		for _, metricDto := range metricsDto {
			if metricDto.MType == metric.CounterType {
				name := metricDto.ID
				value := metric.Counter(*metricDto.Delta)

				if !h.Counter.Verify(ctx, name, value, metricDto.Hash) {
					continue
				}

				counterMetrics = append(counterMetrics, metric.CounterDto{Name: name, Value: value})
			}
		}

		gaugeMetrics := make([]metric.GaugeDto, 0)
		for _, metricDto := range metricsDto {
			if metricDto.MType == metric.GaugeType {
				name := metricDto.ID
				value := metric.Gauge(*metricDto.Value)

				if !h.Gauge.Verify(ctx, name, value, metricDto.Hash) {
					continue
				}

				gaugeMetrics = append(gaugeMetrics, metric.GaugeDto{Name: name, Value: value})
			}
		}

		if len(counterMetrics) != 0 {
			err = h.Counter.UpdateList(ctx, counterMetrics)
			if err != nil {
				h.Log(ctx).Error().Err(err).Msg("invalid save counter metrics")
				http.Error(rw, "invalid save metrics", http.StatusBadRequest)
				return
			}
		}

		if len(gaugeMetrics) != 0 {
			err = h.Gauge.UpdateList(ctx, gaugeMetrics)
			if err != nil {
				h.Log(ctx).Error().Err(err).Msg("invalid save gauge metrics")
				http.Error(rw, "invalid save metrics", http.StatusBadRequest)
				return
			}
		}

		h.Log(ctx).Info().Msg("json update list metrics handled")

		rw.Write([]byte("OK"))
	}
}
