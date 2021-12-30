package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) GetCounterMetricHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "name")

		metricValue, err := h.Counter.GetOne(metricName)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		rw.Write([]byte(strconv.FormatInt(int64(metricValue), 10)))
	}
}

func (h *Handler) GetGaugeMetricHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "name")

		metricValue, err := h.Gauge.GetOne(metricName)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		rw.Write([]byte(strconv.FormatFloat(float64(metricValue), 'f', -1, 64)))
	}
}
