package handler

import (
	"bytes"
	"encoding/json"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/service/mocks"
	commonMocks "github.com/djokcik/praktikum-go-devops/internal/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetJSONHandler(t *testing.T) {
	t.Run("1. Should return error when metric is unknown type", func(t *testing.T) {
		h := Handler{Mux: chi.NewMux()}
		h.Post("/update/", h.GetJSONHandler())

		rBody, _ := json.Marshal(metric.MetricDto{ID: "TestMetric", MType: "TestType"})
		request := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(rBody))

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()

		body, _ := io.ReadAll(res.Body)
		defer res.Body.Close()

		require.Equal(t, string(body), "unknown type metric: TestType\n")
		require.Equal(t, res.StatusCode, http.StatusBadRequest)
	})

	t.Run("2. Should return counter value", func(t *testing.T) {
		m := mocks.CounterService{Mock: mock.Mock{}}
		m.On("GetOne", mock.Anything, "TestMetric").Return(metric.Counter(10), nil)

		mHash := commonMocks.HashService{Mock: mock.Mock{}}
		mHash.On("GetCounterHash", mock.Anything, "TestMetric", metric.Counter(10)).Return("hashTest")

		h := Handler{Counter: &m, Hash: &mHash, Mux: chi.NewMux()}
		h.Post("/update/", h.GetJSONHandler())

		rBody, _ := json.Marshal(metric.MetricDto{ID: "TestMetric", MType: "counter"})
		request := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(rBody))

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()

		body, _ := io.ReadAll(res.Body)
		defer res.Body.Close()

		require.Equal(t, string(body), `{"id":"TestMetric","type":"counter","delta":10,"hash":"hashTest"}`)
		require.Equal(t, res.StatusCode, http.StatusOK)
	})

	t.Run("3. Should update gauge value", func(t *testing.T) {
		m := mocks.GaugeService{Mock: mock.Mock{}}
		m.On("GetOne", mock.Anything, "TestMetric").Return(metric.Gauge(0.123), nil)

		mHash := commonMocks.HashService{Mock: mock.Mock{}}
		mHash.On("GetGaugeHash", mock.Anything, "TestMetric", metric.Gauge(0.123)).Return("hashTest")

		h := Handler{Gauge: &m, Hash: &mHash, Mux: chi.NewMux()}
		h.Post("/update/", h.GetJSONHandler())

		rDto := metric.MetricDto{ID: "TestMetric", MType: "gauge"}
		rBody, _ := json.Marshal(rDto)

		request := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(rBody))

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()

		body, _ := io.ReadAll(res.Body)
		defer res.Body.Close()

		require.Equal(t, string(body), `{"id":"TestMetric","type":"gauge","value":0.123,"hash":"hashTest"}`)
		require.Equal(t, res.StatusCode, http.StatusOK)
	})
}
