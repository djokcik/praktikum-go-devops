package agent

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/google/uuid"
	"net/http"
)

func (a *agent) SendToServer(ctx context.Context) {
	traceID, _ := uuid.NewUUID()
	logger := a.Log(ctx).With().Str(logging.ServiceKey, "SendToServer").Str(logging.TraceIDKey, traceID.String()).Logger()
	ctx = logging.SetCtxLogger(ctx, logger)

	a.Log(ctx).Info().Msg("Start send metrics")

	a.Lock()
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
	a.Unlock()

	url := fmt.Sprintf("http://%s/updates/", a.cfg.Address)
	body, err := json.Marshal(metricDtoList)
	if err != nil {
		a.Log(ctx).Warn().Err(err).Msgf("invalid marshal json: %s", err)
		return
	}

	if a.cfg.PublicKey.PublicKey != nil {
		body, err = a.Hash.EncryptOAEP(sha1.New(), rand.Reader, a.cfg.PublicKey.PublicKey, body, nil)
		if err != nil {
			a.Log(ctx).Error().Err(err).Msgf("encrypt: %s", err)
			return
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		a.Log(ctx).Error().Err(err).Msg("request was interrupted")
		return
	}

	req.Header.Set("Content-Type", "application/json")

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
