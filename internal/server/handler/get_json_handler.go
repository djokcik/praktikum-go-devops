package handler

import (
	"encoding/json"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"net/http"
)

func (h *Handler) GetJSONHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "GetJSONHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		var metricDto metric.MetricsDto
		err := json.NewDecoder(r.Body).Decode(&metricDto)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		logger.UpdateContext(metricDto.GetLoggerContext)
		ctx = logging.SetCtxLogger(ctx, logger)

		switch metricDto.MType {
		default:
			errMessage := fmt.Sprintf("unknown type metric: %s", metricDto.MType)
			h.Log(ctx).Warn().Msgf(errMessage)

			http.Error(rw, errMessage, http.StatusBadRequest)
			return
		case metric.CounterType:
			val, _ := h.Counter.GetOne(ctx, metricDto.ID)
			delta := int64(val)

			metricDto.Delta = &delta
		case metric.GaugeType:
			val, _ := h.Gauge.GetOne(ctx, metricDto.ID)
			value := float64(val)

			metricDto.Value = &value
		}

		rw.Header().Set("Content-Type", "application/json")

		h.Log(ctx).Info().Msg("json metrics handled")

		bytes, _ := json.Marshal(metricDto)
		rw.Write(bytes)
	}
}
