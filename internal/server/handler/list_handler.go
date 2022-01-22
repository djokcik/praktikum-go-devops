package handler

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"html/template"
	"net/http"
)

func (h *Handler) ListHandler() http.HandlerFunc {
	type listTemplateData struct {
		CounterMetrics []metric.Metric
		GaugeMetrics   []metric.Metric
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "ListHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		counterMetrics, err := h.Counter.List(ctx)
		if err != nil {
			h.Log(ctx).Error().Err(err).Msg("counterMetrics internal error")
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		gaugeMetrics, err := h.Gauge.List(ctx)
		if err != nil {
			h.Log(ctx).Error().Err(err).Msg("gaugeMetrics internal error")
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("static/listMetrics.html")
		if err != nil {
			h.Log(ctx).Error().Err(err).Msg("read template internal error")
			http.Error(rw, "Internal Server Error", 500)
			return
		}

		rw.Header().Set("Content-Type", "text/html")

		err = tmpl.Execute(rw, listTemplateData{
			GaugeMetrics:   gaugeMetrics,
			CounterMetrics: counterMetrics,
		})

		if err != nil {
			h.Log(ctx).Error().Err(err).Msg("parse template internal error")
			http.Error(rw, "Internal Server Error", 500)
		}

		h.Log(ctx).Info().Msg("list handled")
	}
}
