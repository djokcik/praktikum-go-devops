package handler

import (
	"bytes"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_UpdateListJSONHandler(t *testing.T) {
	t.Run("1. Should update counter and gauge list", func(t *testing.T) {
		mCounter := mocks.CounterService{Mock: mock.Mock{}}
		mCounter.On("Verify", mock.Anything, "TestCounterMetric", metric.Counter(10), "hashCounter").
			Return(true)
		mCounter.On("UpdateList", mock.Anything, []metric.CounterDto{{Name: "TestCounterMetric", Value: metric.Counter(10)}}).
			Return(nil)

		mGauge := mocks.GaugeService{Mock: mock.Mock{}}
		mGauge.On("Verify", mock.Anything, "TestGaugeMetric", metric.Gauge(1.5), "hashGauge").
			Return(true)
		mGauge.On("UpdateList", mock.Anything, []metric.GaugeDto{{Name: "TestGaugeMetric", Value: metric.Gauge(1.5)}}).
			Return(nil)

		h := Handler{Mux: chi.NewMux(), Counter: &mCounter, Gauge: &mGauge}
		h.Post("/updates/", h.UpdateListJSONHandler())

		requestBody := `[{"id":"TestCounterMetric","type":"counter","delta":10,"hash":"hashCounter"},{"id":"TestGaugeMetric","type":"gauge","value":1.5,"hash":"hashGauge"}]`
		request := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewBufferString(requestBody))

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusOK)
		mCounter.AssertNumberOfCalls(t, "Verify", 1)
		mCounter.AssertNumberOfCalls(t, "UpdateList", 1)
		mGauge.AssertNumberOfCalls(t, "Verify", 1)
		mGauge.AssertNumberOfCalls(t, "UpdateList", 1)
	})
}
