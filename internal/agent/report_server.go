package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"log"
	"net/http"
)

func (a *agent) SendToServer(ctx context.Context) func() {
	return func() {
		for _, sendMetric := range a.CollectedMetric {
			if ctx.Err() != nil {
				log.Printf("context cancelled")
				break
			}

			metricName := sendMetric.Name
			metricType := sendMetric.Type
			metricValue := sendMetric.Value

			url := fmt.Sprintf("http://%s/update/", a.cfg.Address)

			var metricDto metric.MetricsDto
			switch metricType {
			default:
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
				log.Printf("request was interrapted with error: %s", err)
				continue
			}

			req.Header.Set("Content-Type", "application/json")

			res, err := a.Client.Do(req)
			if err != nil {
				log.Printf("request %s ended with error %v", url, err)
				continue
			}

			err = res.Body.Close()
			if err != nil {
				log.Printf("read from body closed with error: %v", err)
			}
		}
	}
}
