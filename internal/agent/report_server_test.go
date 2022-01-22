package agent

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendToServer(t *testing.T) {
	t.Run("Should send metrics to server", func(t *testing.T) {
		metricAgent := NewAgent(Config{Address: "127.0.0.1:45555"})

		collectedMap := make(map[string]SendAgentMetric)
		collectedMap["TestMetric"] = SendAgentMetric{Name: "TestMetric", Type: "counter", Value: metric.Counter(10)}

		metricAgent.CollectedMetric = collectedMap

		l, err := net.Listen("tcp", "127.0.0.1:45555")
		if err != nil {
			log.Fatal(err)
		}

		// Start a local HTTP server
		ts := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request parameters
			defer req.Body.Close()
			body, _ := io.ReadAll(req.Body)

			require.Equal(t, string(body), `{"id":"TestMetric","type":"counter","delta":10}`)
			require.Equal(t, req.Method, http.MethodPost)
			require.Equal(t, req.URL.String(), "/update/")
			require.Equal(t, req.Header.Get("Content-Type"), "application/json")
			// Send response to be tested
			rw.Write([]byte(`OK`))
		}))

		ts.Listener.Close()
		ts.Listener = l

		// Start the server.
		ts.Start()

		// Close the server when test finishes
		defer ts.Close()

		metricAgent.SendToServer(context.Background())()
	})
}
