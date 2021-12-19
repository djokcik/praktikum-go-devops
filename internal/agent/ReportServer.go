package agent

import (
	"fmt"
	"net/http"
	"time"
)

const reportInterval = 10
const host = "127.0.0.1"
const port = "8080"

func ReportMetricsToServer(updatedMetric map[string]SendAgentMetric) {
	ticker := time.NewTicker(reportInterval * time.Second)
	client := http.Client{}

	for range ticker.C {
		SendToServer(updatedMetric, client)
	}
}

func SendToServer(updatedMetric map[string]SendAgentMetric, client http.Client) {
	for _, sendMetric := range updatedMetric {
		metric := sendMetric.Metric
		value := sendMetric.Value

		url := fmt.Sprintf("http://%s:%s/update/%s/%s/%v", host, port, metric.Type(), metric.Name(), value)

		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			fmt.Printf("Запрос был прерван с ошибков: %s", err)
			continue
		}

		req.Header.Set("application-type", "text/plain")

		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("Запрос %s завершился с ошибкой %v", url, err)
			continue
		}

		err = res.Body.Close()
		if err != nil {
			fmt.Printf("Чтение из body закрылось с ошибкой: %v", err)
		}
	}
}
