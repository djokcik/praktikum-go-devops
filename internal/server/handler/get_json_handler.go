package handler

import (
	"encoding/json"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"net/http"
)

func (h *Handler) GetJSONHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var metricDto metric.MetricsDto
		err := decoder.Decode(&metricDto)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		switch metricDto.MType {
		default:
			http.Error(rw, fmt.Sprintf("unknown type metric: %s", metricDto.MType), http.StatusBadRequest)
			return
		case metric.CounterType:
			val, _ := h.Counter.GetOne(metricDto.ID)
			delta := int64(val)

			metricDto.Delta = &delta
		case metric.GaugeType:
			val, _ := h.Gauge.GetOne(metricDto.ID)
			value := float64(val)

			metricDto.Value = &value
		}

		rw.Header().Set("Content-Type", "application/json")

		bytes, _ := json.Marshal(metricDto)
		rw.Write(bytes)
	}
}
