package handlers

import (
	"fmt"
	"github.com/Jokcik/praktikum-go-devops/internal/metrics"
	"github.com/Jokcik/praktikum-go-devops/internal/server/servermetrics"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage"
	"net/http"
	"strconv"
)

func CounterHandler(repository storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		counterMetric := servermetrics.GetServerCounterMetrics()

		parseAndSaveValueFunc := func(name string, value string) {
			parseValue, err := strconv.ParseInt(value, 10, 64)

			if err != nil {
				fmt.Printf("Error %v", err)
				http.Error(rw, "invalid value", http.StatusBadRequest)
			}

			repository.Update(name, metrics.Counter(parseValue))
		}

		MetricHandler(counterMetric, parseAndSaveValueFunc)(rw, r)
	}
}
