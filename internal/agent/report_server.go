package agent

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	pb "github.com/djokcik/praktikum-go-devops/pkg/proto"
	"github.com/google/uuid"
	"net"
	"net/http"
	"os"
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

	if a.GRPCClient != nil {
		var metrics []*pb.Metric

		for _, metricDto := range metricDtoList {
			pbMetric := pb.Metric{
				ID:    metricDto.ID,
				Hash:  metricDto.Hash,
				MType: metricDto.MType,
			}

			if metricDto.Delta != nil {
				pbMetric.Delta = *metricDto.Delta
			}

			if metricDto.Value != nil {
				pbMetric.Value = *metricDto.Value
			}

			metrics = append(metrics, &pbMetric)
		}

		res, err := a.GRPCClient.SendMetric(ctx, &pb.SendMetricRequest{Metrics: metrics})
		if err != nil {
			a.Log(ctx).Error().Err(err).Msg("gRPC request ended with error")
			return
		}

		if res.Error != "" {
			a.Log(ctx).Error().Err(errors.New(res.Error)).Msg("gRPC request ended with error")
			return
		}
	} else {
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
		req.Header.Set("X-Real-IP", GetLocalIP().String())

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
	}

	a.Log(ctx).Info().Msg("Finished send metrics")
}

func GetLocalIP() net.IP {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4
		}
	}

	return nil
}
