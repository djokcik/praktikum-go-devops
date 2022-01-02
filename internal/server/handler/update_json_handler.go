package handler

import (
	"encoding/json"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"net/http"
)

func (h *Handler) UpdateJSONHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var metricDto metric.MetricsDto
		err := json.NewDecoder(r.Body).Decode(&metricDto)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		switch metricDto.MType {
		default:
			http.Error(rw, fmt.Sprintf("unknown type metric: %s", metricDto.MType), http.StatusBadRequest)
			return
		case metric.CounterType:
			err = h.Counter.Increase(metricDto.ID, metric.Counter(*metricDto.Delta))
		case metric.GaugeType:
			_, err = h.Gauge.Update(metricDto.ID, metric.Gauge(*metricDto.Value))
		}

		if err != nil {
			// TODO: add log
			http.Error(rw, "invalid save metric", http.StatusBadRequest)
			return
		}

		rw.Write([]byte("OK"))
	}
}
