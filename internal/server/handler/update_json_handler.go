package handler

import (
	"encoding/json"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"net/http"
)

func (h *Handler) UpdateJSONHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "UpdateJSONHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		var metricDto metric.MetricDto
		err := json.NewDecoder(r.Body).Decode(&metricDto)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, logger = logging.GetCtxLogger(ctx)
		logger.UpdateContext(metricDto.GetLoggerContext)
		ctx = logging.SetCtxLogger(ctx, logger)

		switch metricDto.MType {
		default:
			errMessage := fmt.Sprintf("unknown type metric: %s", metricDto.MType)
			h.Log(ctx).Warn().Msg(errMessage)

			http.Error(rw, errMessage, http.StatusBadRequest)
			return
		case metric.CounterType:
			name := metricDto.ID
			value := metric.Counter(*metricDto.Delta)

			if !h.Counter.Verify(ctx, name, value, metricDto.Hash) {
				http.Error(rw, "Invalid sign", http.StatusBadRequest)
				return
			}

			err = h.Counter.Increase(ctx, name, value)
		case metric.GaugeType:
			name := metricDto.ID
			value := metric.Gauge(*metricDto.Value)

			if !h.Gauge.Verify(ctx, name, value, metricDto.Hash) {
				http.Error(rw, "Invalid sign", http.StatusBadRequest)
				return
			}

			err = h.Gauge.Update(ctx, name, value)
		}

		if err != nil {
			h.Log(ctx).Error().Err(err).Msg("invalid save metric")
			http.Error(rw, "invalid save metric", http.StatusBadRequest)
			return
		}

		h.Log(ctx).Info().Msg("json metrics handled")

		rw.Write([]byte("OK"))
	}
}
