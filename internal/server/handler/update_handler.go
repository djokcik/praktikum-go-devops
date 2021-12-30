package handler

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) GaugeHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "name")
		metricValue := chi.URLParam(r, "value")

		parseValue, err := strconv.ParseFloat(metricValue, 64)

		if err != nil {
			http.Error(rw, "invalid value", http.StatusBadRequest)
			return
		}

		_, err = h.Gauge.Update(metricName, metric.Gauge(parseValue))
		if err != nil {
			http.Error(rw, "invalid save metric", http.StatusBadRequest)
			return
		}

		rw.Write([]byte("OK"))
	}
}

func (h *Handler) CounterHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "name")
		metricValue := chi.URLParam(r, "value")

		parseValue, err := strconv.ParseInt(metricValue, 0, 64)

		if err != nil {
			http.Error(rw, "invalid value", http.StatusBadRequest)
			return
		}

		err = h.Counter.AddValue(metricName, metric.Counter(parseValue))

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		rw.Write([]byte("OK"))
	}
}
