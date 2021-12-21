package handler

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) CounterHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "name")
		metricValue := chi.URLParam(r, "value")

		parseValue, err := strconv.ParseInt(metricValue, 10, 64)

		if err != nil {
			http.Error(rw, "invalid value", http.StatusBadRequest)
			rw.Write([]byte("invalid value"))
			return
		}

		_, err = h.Repo.Update(metricType, metric.Counter(parseValue))
		if err != nil {
			http.Error(rw, "invalid save metric", http.StatusBadRequest)
			rw.Write([]byte("invalid save metric"))
			return
		}

		rw.Write([]byte("OK"))
	}
}
