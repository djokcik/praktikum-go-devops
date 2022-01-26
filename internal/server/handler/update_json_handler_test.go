package handler

import (
	"bytes"
	"encoding/json"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_UpdateJSONHandler(t *testing.T) {
	t.Run("1. Should return error when json is invalid", func(t *testing.T) {
		h := Handler{Mux: chi.NewMux()}
		h.Post("/update/", h.UpdateJSONHandler())

		requestBody := `{"ID":"TestMetric","MType":"Counter","Delta":"10}`
		request := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBufferString(requestBody))

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusBadRequest)
	})

	t.Run("2. Should return error when metric is unknown type", func(t *testing.T) {
		h := Handler{Mux: chi.NewMux()}
		h.Post("/update/", h.UpdateJSONHandler())

		delta := int64(10)
		rDto := metric.MetricDto{ID: "TestMetric", MType: "TestType", Delta: &delta}
		rBody, _ := json.Marshal(rDto)

		request := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(rBody))

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()

		body, _ := io.ReadAll(res.Body)
		defer res.Body.Close()

		require.Equal(t, string(body), "unknown type metric: TestType\n")
		require.Equal(t, res.StatusCode, http.StatusBadRequest)
	})

	t.Run("3. Should update counter value", func(t *testing.T) {
		m := mocks.CounterService{Mock: mock.Mock{}}
		m.On("Increase", mock.Anything, "TestMetric", metric.Counter(10)).Return(nil)
		m.On("Verify", mock.Anything, "TestMetric", metric.Counter(10), "myHash").Return(true)

		h := Handler{Counter: &m, Mux: chi.NewMux()}
		h.Post("/update/", h.UpdateJSONHandler())

		delta := int64(10)
		rDto := metric.MetricDto{ID: "TestMetric", MType: "counter", Delta: &delta, Hash: "myHash"}
		rBody, _ := json.Marshal(rDto)

		request := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(rBody))

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusOK)
		m.AssertNumberOfCalls(t, "Increase", 1)
		m.AssertNumberOfCalls(t, "Verify", 1)
	})

	t.Run("4. Should update gauge value", func(t *testing.T) {
		m := mocks.GaugeService{Mock: mock.Mock{}}
		m.On("Update", mock.Anything, "TestMetric", metric.Gauge(0.123)).Return(nil)
		m.On("Verify", mock.Anything, "TestMetric", metric.Gauge(0.123), "myHash").Return(true)

		h := Handler{Gauge: &m, Mux: chi.NewMux()}
		h.Post("/update/", h.UpdateJSONHandler())

		delta := 0.123
		rDto := metric.MetricDto{ID: "TestMetric", MType: "gauge", Value: &delta, Hash: "myHash"}
		rBody, _ := json.Marshal(rDto)

		request := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(rBody))

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusOK)
		m.AssertNumberOfCalls(t, "Update", 1)
		m.AssertNumberOfCalls(t, "Verify", 1)
	})
}
