package handler

import (
	"fmt"
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) GetGaugeHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "name")

		metricValue, err := h.Repo.Get(metricType)
		fmt.Println("GAUGE", metricType, metricValue)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		switch elem := metricValue.(type) {
		default:
			http.Error(rw, fmt.Sprintf("the metric `%v` didn`t find", metricType), http.StatusNotFound)
			return
		case metric.Gauge:
			rw.Write([]byte(fmt.Sprintf("%v", elem)))
		}
	}
}
