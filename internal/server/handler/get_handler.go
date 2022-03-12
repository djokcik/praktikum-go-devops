package handler

import (
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// GetCounterMetricHandler handler return counter metric
func (h *Handler) GetCounterMetricHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "GetCounterMetricHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		metricName := chi.URLParam(r, "name")

		metricValue, err := h.Counter.GetOne(ctx, metricName)

		logger.UpdateContext(metricValue.GetLoggerContext(metricName))
		ctx = logging.SetCtxLogger(ctx, logger)

		if err != nil {
			h.Log(ctx).Warn().Err(err).Msg("metric not found")
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		h.Log(ctx).Info().Msg("metric handled")

		rw.Write([]byte(strconv.FormatInt(int64(metricValue), 10)))
	}
}

// GetGaugeMetricHandler handler return gauge metric
func (h *Handler) GetGaugeMetricHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "GetGaugeMetricHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		metricName := chi.URLParam(r, "name")
		metricValue, err := h.Gauge.GetOne(ctx, metricName)

		logger.UpdateContext(metricValue.GetLoggerContext(metricName))
		ctx = logging.SetCtxLogger(ctx, logger)

		if err != nil {
			h.Log(ctx).Warn().Err(err).Msg("metric not found")
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		h.Log(ctx).Info().Msg("metric handled")

		rw.Write([]byte(strconv.FormatFloat(float64(metricValue), 'f', -1, 64)))
	}
}
