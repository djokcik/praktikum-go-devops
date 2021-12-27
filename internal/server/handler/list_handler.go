package handler

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"html/template"
	"log"
	"net/http"
)

func (h *Handler) ListHandler() http.HandlerFunc {
	type listTemplateData struct {
		CounterMetrics []metric.Metric
		GaugeMetrics   []metric.Metric
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		counterMetrics, err := h.Repo.List(storage.ListRepositoryFilter{Type: metric.CounterType})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		gaugeMetrics, err := h.Repo.List(storage.ListRepositoryFilter{Type: metric.GaugeType})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("static/listMetrics.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(rw, "Internal Server Error", 500)
			return
		}

		err = tmpl.Execute(rw, listTemplateData{
			GaugeMetrics:   gaugeMetrics.([]metric.Metric),
			CounterMetrics: counterMetrics.([]metric.Metric),
		})

		if err != nil {
			log.Println(err.Error())
			http.Error(rw, "Internal Server Error", 500)
		}
	}
}
