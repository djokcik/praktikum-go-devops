package handler

import (
	"fmt"
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) GaugeHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "name")
		metricValue := chi.URLParam(r, "value")

		parseValue, err := strconv.ParseFloat(metricValue, 64)

		fmt.Println("SAVE GAUGE", metricValue, metricType)
		if err != nil {
			http.Error(rw, "invalid value", http.StatusBadRequest)
			return
		}

		_, err = h.Repo.Update(metricType, metric.Gauge(parseValue))
		if err != nil {
			http.Error(rw, "invalid save metric", http.StatusBadRequest)
			return
		}

		rw.Write([]byte("OK"))
	}
}
