package handler

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) CounterHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "name")
		metricValue := chi.URLParam(r, "value")

		parseValue, err := strconv.ParseInt(metricValue, 0, 64)

		if err != nil {
			http.Error(rw, "invalid value", http.StatusBadRequest)
			rw.Write([]byte("invalid value"))
			return
		}

		val, _ := h.Repo.Get(storage.GetRepositoryFilter{
			Name:         metricName,
			Type:         metric.CounterType,
			DefaultValue: metric.Counter(0),
		})

		_, err = h.Repo.Update(metricName, val.(metric.Counter)+metric.Counter(parseValue))
		if err != nil {
			http.Error(rw, "invalid save metric", http.StatusBadRequest)
			rw.Write([]byte("invalid save metric"))
			return
		}

		rw.Write([]byte("OK"))
	}
}
