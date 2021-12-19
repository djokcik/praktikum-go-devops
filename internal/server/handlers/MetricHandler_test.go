package handlers

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metrics/mocks"
	"github.com/Jokcik/praktikum-go-devops/internal/server/servermetrics"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricHandler(t *testing.T) {
	t.Run("1. Should return error when path is invalid", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/update/counter", nil)

		w := httptest.NewRecorder()
		h := MetricHandler(make(servermetrics.ServerMapMetrics), func(_ string, _ string) {})

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusNotFound)
	})

	t.Run("2. Should return error when method is Get", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/update/counter", nil)

		w := httptest.NewRecorder()
		h := MetricHandler(make(servermetrics.ServerMapMetrics), func(_ string, _ string) {})

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusNotFound)
	})

	t.Run("3. Should be called function callback", func(t *testing.T) {
		m := mocks.Metric{Mock: mock.Mock{}}
		m.On("CallbackFunc").Return(nil)

		serverMap := make(servermetrics.ServerMapMetrics)
		serverMap["Alloc"] = &m

		request := httptest.NewRequest(http.MethodPost, "/update/counter/Alloc/0.123", nil)

		w := httptest.NewRecorder()
		h := MetricHandler(serverMap, func(_ string, _ string) {
			m.MethodCalled("CallbackFunc")
		})

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusOK)
		m.AssertNumberOfCalls(t, "CallbackFunc", 1)
	})
}
