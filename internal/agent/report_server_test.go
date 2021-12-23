package agent

import (
	"context"
	"github.com/stretchr/testify/require"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendToServer(t *testing.T) {
	t.Run("Should send metrics to server", func(t *testing.T) {
		metricAgent := NewAgent(context.Background())

		collectedMap := make(map[string]SendAgentMetric)
		collectedMap["TestMetric"] = SendAgentMetric{Name: "TestMetric", Type: "TestType", Value: "TestValue"}

		metricAgent.CollectedMetric = collectedMap

		l, err := net.Listen("tcp", "127.0.0.1:8080")
		if err != nil {
			log.Fatal(err)
		}

		// Start a local HTTP server
		ts := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request parameters
			require.Equal(t, req.Method, http.MethodPost)
			require.Equal(t, req.URL.String(), "/update/TestType/TestMetric/TestValue")
			require.Equal(t, req.Header.Get("application-type"), "text/plain")
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
