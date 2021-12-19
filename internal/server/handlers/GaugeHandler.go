package handlers

import (
	"fmt"
	"github.com/Jokcik/praktikum-go-devops/internal/metrics"
	"github.com/Jokcik/praktikum-go-devops/internal/server/servermetrics"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage"
	"net/http"
	"strconv"
)

func GaugeHandler(repository storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		gaugeMetric := servermetrics.GetServerGaugeMetrics()

		parseAndSaveValueFunc := func(name string, value string) {
			parseValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				fmt.Printf("Error %v", err)
				http.Error(rw, "invalid value", http.StatusBadRequest)
				return
			}

			_, err = repository.Update(name, metrics.Gauge(parseValue))
			if err != nil {
				http.Error(rw, "invalid save metrics", http.StatusBadRequest)
			}
		}

		MetricHandler(gaugeMetric, parseAndSaveValueFunc)(rw, r)
	}
}
