package handler

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
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
		counterMetrics, err := h.Counter.List()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		gaugeMetrics, err := h.Gauge.List()
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
			GaugeMetrics:   gaugeMetrics,
			CounterMetrics: counterMetrics,
		})

		if err != nil {
			log.Println(err.Error())
			http.Error(rw, "Internal Server Error", 500)
		}
	}
}
