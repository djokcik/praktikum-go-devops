package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/google/uuid"
	"net/http"
)

func (a *agent) SendToServer(ctx context.Context) func() {
	return func() {
		traceID, _ := uuid.NewUUID()
		logger := a.Log(ctx).With().Str(logging.ServiceKey, "SendToServer").Str(logging.TraceIDKey, traceID.String()).Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		a.Log(ctx).Info().Msg("Start send metrics")
		for _, sendMetric := range a.CollectedMetric {
			if ctx.Err() != nil {
				a.Log(ctx).Error().Err(ctx.Err())
				break
			}

			metricName := sendMetric.Name
			metricType := sendMetric.Type
			metricValue := sendMetric.Value

			url := fmt.Sprintf("http://%s/update/", a.cfg.Address)

			var metricDto metric.MetricsDto
			switch metricType {
			default:
				a.Log(ctx).Error().Msgf("Invalid metric type: %s", metricType)
				continue
			case metric.GaugeType:
				value := float64(metricValue.(metric.Gauge))
				metricDto = metric.MetricsDto{ID: metricName, MType: metricType, Value: &value}
			case metric.CounterType:
				delta := int64(metricValue.(metric.Counter))
				metricDto = metric.MetricsDto{ID: metricName, MType: metricType, Delta: &delta}
			}

			body, _ := json.Marshal(metricDto)
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
			if err != nil {
				a.Log(ctx).Error().Err(err).Msg("request was interrupted")
				continue
			}

			req.Header.Set("Content-Type", "application/json")

			res, err := a.Client.Do(req)
			if err != nil {
				a.Log(ctx).Error().Err(err).Msg("request ended with error")
				continue
			}

			err = res.Body.Close()
			if err != nil {
				a.Log(ctx).Error().Err(err).Msg("read from body closed with error")
				continue
			}
		}
		a.Log(ctx).Info().Msg("Finished send metrics")
	}
}
