package agent

import (
	"bytes"
	"compress/gzip"
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

		var metricDtoList []metric.MetricDto
		for _, sendMetric := range a.CollectedMetric {
			if ctx.Err() != nil {
				a.Log(ctx).Error().Err(ctx.Err())
				break
			}

			metricName := sendMetric.Name
			metricType := sendMetric.Type
			metricValue := sendMetric.Value

			var metricDto metric.MetricDto
			switch metricType {
			default:
				a.Log(ctx).Error().Msgf("Invalid metric type: %s", metricType)
				continue
			case metric.GaugeType:
				value := metricValue.(metric.Gauge)
				refValue := float64(value)
				hash := a.Hash.GetGaugeHash(ctx, metricName, value)
				metricDto = metric.MetricDto{ID: metricName, MType: metricType, Value: &refValue, Hash: hash}
			case metric.CounterType:
				delta := metricValue.(metric.Counter)
				refDelta := int64(metricValue.(metric.Counter))
				hash := a.Hash.GetCounterHash(ctx, metricName, delta)
				metricDto = metric.MetricDto{ID: metricName, MType: metricType, Delta: &refDelta, Hash: hash}
			}

			metricDtoList = append(metricDtoList, metricDto)
		}

		url := fmt.Sprintf("http://%s/updates/", a.cfg.Address)
		body, _ := json.Marshal(metricDtoList)
		var buf bytes.Buffer
		g := gzip.NewWriter(&buf)

		if _, err := g.Write(body); err != nil {
			a.Log(ctx).Error().Err(err).Msg("Invalid Write gzip")
			return
		}

		if err := g.Close(); err != nil {
			a.Log(ctx).Error().Err(err).Msg("Invalid close gzip")
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
		if err != nil {
			a.Log(ctx).Error().Err(err).Msg("request was interrupted")
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")

		res, err := a.Client.Do(req)
		if err != nil {
			a.Log(ctx).Error().Err(err).Msg("request ended with error")
			return
		}

		err = res.Body.Close()
		if err != nil {
			a.Log(ctx).Error().Err(err).Msg("read from body closed with error")
			return
		}

		a.Log(ctx).Info().Msg("Finished send metrics")
	}
}
