package handler

import (
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) GetMetricHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")

		metricValue, err := h.Repo.Get(storage.GetRepositoryFilter{Type: metricType, Name: metricName})

		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		switch metricType {
		default:
			http.Error(rw, fmt.Sprintf("the metric `%v` didn`t find", metricType), http.StatusNotFound)
		case metric.GaugeType:
			rw.Write([]byte(strconv.FormatFloat(float64(metricValue.(metric.Gauge)), 'f', -1, 64)))
		case metric.CounterType:
			rw.Write([]byte(strconv.FormatInt(int64(metricValue.(metric.Counter)), 10)))
		}
	}
}
