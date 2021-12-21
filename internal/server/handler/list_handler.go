package handler

import (
	"fmt"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage"
	"net/http"
)

func (h *Handler) ListHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		list, err := h.Repo.List()

		if err != nil {
			http.Error(rw, "invalid get metrics", http.StatusBadGateway)
			return
		}

		result := "<html><body>"
		metricList := list.([]storage.MetricElement)

		if len(metricList) != 0 {
			result += "<ul>"
			for _, metric := range list.([]storage.MetricElement) {
				result += fmt.Sprintf("<li>Name: %v<br>Value: %v</li>", metric.Name, metric.Value)
			}
			result += "</ul>"
		} else {
			result += "<h1>Нет метрик</h1>"
		}

		result += "</body></html>"

		rw.Write([]byte(result))
	}
}
