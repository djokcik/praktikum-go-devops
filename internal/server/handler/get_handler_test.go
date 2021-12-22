package handler

import (
	"errors"
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetMetricHandler(t *testing.T) {
	t.Run("1. Should return gauge metric", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Get", storage.GetRepositoryFilter{Type: metric.GaugeType, Name: "TestName"}).Return(metric.Gauge(0.123), nil)

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodGet, "/value/gauge/TestName", nil)
		h.Get("/value/{type}/{name}", h.GetMetricHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		m.AssertNumberOfCalls(t, "Get", 1)
		require.Equal(t, res.StatusCode, http.StatusOK)
		require.Equal(t, string(resBody), "0.123")
	})

	t.Run("2. Should return counter metric", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Get", storage.GetRepositoryFilter{Type: metric.CounterType, Name: "TestName"}).Return(metric.Counter(123), nil)

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodGet, "/value/counter/TestName", nil)
		h.Get("/value/{type}/{name}", h.GetMetricHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		m.AssertNumberOfCalls(t, "Get", 1)
		require.Equal(t, res.StatusCode, http.StatusOK)
		require.Equal(t, string(resBody), "123")
	})

	t.Run("3. Should return 404 when metric didn`t find", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Get", storage.GetRepositoryFilter{Type: metric.GaugeType, Name: "TestName"}).Return(nil, errors.New("error"))

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodGet, "/value/gauge/TestName", nil)
		h.Get("/value/{type}/{name}", h.GetMetricHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		m.AssertNumberOfCalls(t, "Get", 1)
		require.Equal(t, res.StatusCode, http.StatusNotFound)
	})
}
