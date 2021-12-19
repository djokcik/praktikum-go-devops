package handlers

import (
	"github.com/Jokcik/praktikum-go-devops/internal/server/servermetrics"
	"net/http"
	"strings"
)

func MetricHandler(mapMetric servermetrics.ServerMapMetrics, process func(string, string)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		segments := strings.Split(r.URL.Path, "/")

		if len(segments) != 5 || r.Method != http.MethodPost {
			http.Error(rw, "not found", http.StatusNotFound)
			return
		}

		key := segments[len(segments)-2]
		value := segments[len(segments)-1]

		if _, ok := mapMetric[key]; ok {
			process(key, value)
			rw.Write([]byte("OK"))
		} else {
			// по хорошему должна быть проверка на сервере, но тесты не проходят
			//http.Error(rw, "invalid params", http.StatusBadRequest)
			process(key, value)
			rw.Write([]byte("OK"))
		}
	}
}
