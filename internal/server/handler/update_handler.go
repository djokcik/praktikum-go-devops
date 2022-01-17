package handler

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) GaugeHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "GaugeHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		metricName := chi.URLParam(r, "name")
		metricValue := chi.URLParam(r, "value")

		parseValue, err := strconv.ParseFloat(metricValue, 64)

		if err != nil {
			h.Log(ctx).Warn().Err(err).Msgf("invalid metricValue: %s", metricValue)
			http.Error(rw, "invalid value", http.StatusBadRequest)
			return
		}

		_, err = h.Gauge.Update(ctx, metricName, metric.Gauge(parseValue))
		if err != nil {
			h.Log(ctx).Err(err).Msgf("invalid save metricName: %s and metricValue: %s", metricName, metricValue)
			http.Error(rw, "invalid save metric", http.StatusBadRequest)
			return
		}

		h.Log(ctx).Info().Msg("update handled")

		rw.Write([]byte("OK"))
	}
}

func (h *Handler) CounterHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "CounterHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		metricName := chi.URLParam(r, "name")
		metricValue := chi.URLParam(r, "value")

		parseValue, err := strconv.ParseInt(metricValue, 0, 64)

		if err != nil {
			h.Log(ctx).Warn().Err(err).Msgf("invalid metricValue: %s", metricValue)
			http.Error(rw, "invalid value", http.StatusBadRequest)
			return
		}

		err = h.Counter.Increase(ctx, metricName, metric.Counter(parseValue))

		if err != nil {
			h.Log(ctx).Warn().Err(err).Msgf("invalid save metricName: %s and metricValue: %s", metricName, metricValue)
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		h.Log(ctx).Info().Msg("update handled")

		rw.Write([]byte("OK"))
	}
}
