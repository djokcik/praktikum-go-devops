package handler

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/service/mocks"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetCounterMetricHandler(t *testing.T) {
	t.Run("1. Should return 404 when metric didn`t find", func(t *testing.T) {
		m := mocks.CounterService{Mock: mock.Mock{}}
		m.On("GetOne", mock.Anything, "TestName").Return(metric.Counter(0), storage.ErrValueNotFound)

		h := Handler{Counter: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodGet, "/value/counter/TestName", nil)
		h.Get("/value/counter/{name}", h.GetCounterMetricHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		m.AssertNumberOfCalls(t, "GetOne", 1)
		require.Equal(t, res.StatusCode, http.StatusNotFound)
	})

	t.Run("2. Should return 404 when metric didn`t find", func(t *testing.T) {
		m := mocks.CounterService{Mock: mock.Mock{}}
		m.On("GetOne", mock.Anything, "TestName").Return(metric.Counter(0), storage.ErrValueNotFound)

		h := Handler{Counter: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodGet, "/value/counter/TestName", nil)
		h.Get("/value/counter/{name}", h.GetCounterMetricHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		m.AssertNumberOfCalls(t, "GetOne", 1)
		require.Equal(t, res.StatusCode, http.StatusNotFound)
	})
}

func TestHandler_GetGaugeMetricHandler(t *testing.T) {
	t.Run("1. Should return gauge metric", func(t *testing.T) {
		m := mocks.GaugeService{Mock: mock.Mock{}}
		m.On("GetOne", mock.Anything, "TestName").Return(metric.Gauge(0.123), nil)

		h := Handler{Gauge: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodGet, "/value/gauge/TestName", nil)
		h.Get("/value/gauge/{name}", h.GetGaugeMetricHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		m.AssertNumberOfCalls(t, "GetOne", 1)
		require.Equal(t, res.StatusCode, http.StatusOK)
		require.Equal(t, string(resBody), "0.123")
	})

	t.Run("2. Should return 404 when metric didn`t find", func(t *testing.T) {
		m := mocks.GaugeService{Mock: mock.Mock{}}
		m.On("GetOne", mock.Anything, "TestName").Return(metric.Gauge(0), storage.ErrValueNotFound)

		h := Handler{Gauge: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodGet, "/value/gauge/TestName", nil)
		h.Get("/value/gauge/{name}", h.GetGaugeMetricHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		m.AssertNumberOfCalls(t, "GetOne", 1)
		require.Equal(t, res.StatusCode, http.StatusNotFound)
	})
}
