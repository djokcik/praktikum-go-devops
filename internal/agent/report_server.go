package agent

import (
	"fmt"
	"log"
	"net/http"
)

const host = "127.0.0.1"
const port = "8080"

func (a *agent) SendToServer() {
	for _, sendMetric := range a.CollectedMetric {
		if a.Ctx.Err() != nil {
			log.Printf("context cancelled")
			break
		}

		metricName := sendMetric.Name
		metricType := sendMetric.Type
		metricValue := sendMetric.Value

		url := fmt.Sprintf("http://%s:%s/update/%s/%s/%v", host, port, metricType, metricName, metricValue)

		req, err := http.NewRequestWithContext(a.Ctx, http.MethodPost, url, nil)
		if err != nil {
			log.Printf("request was interrapted with error: %s", err)
			continue
		}

		req.Header.Set("application-type", "text/plain")

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
