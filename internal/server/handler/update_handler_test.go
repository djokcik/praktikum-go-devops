package handler

import (
	"errors"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GaugeHandler(t *testing.T) {
	t.Run("1. Should update metric", func(t *testing.T) {
		m := mocks.GaugeService{Mock: mock.Mock{}}
		m.On("Update", "Alloc", metric.Gauge(0.123)).Return(true, nil)

		h := Handler{Gauge: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/0.123", nil)
		h.Post("/update/gauge/{name}/{value}", h.GaugeHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusOK)
		m.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("2. Should return error when value is string", func(t *testing.T) {
		m := mocks.GaugeService{Mock: mock.Mock{}}
		m.On("Update", mock.Anything, mock.Anything).Return(true, nil)

		h := Handler{Gauge: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/test", nil)
		h.Post("/update/gauge/{name}/{value}", h.GaugeHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusBadRequest)
		m.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("3. Should return error when update was error", func(t *testing.T) {
		m := mocks.GaugeService{Mock: mock.Mock{}}
		m.On("Update", mock.Anything, mock.Anything).Return(false, errors.New("error"))

		h := Handler{Gauge: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/0.123", nil)
		h.Post("/update/gauge/{name}/{value}", h.GaugeHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusBadRequest)
		m.AssertNumberOfCalls(t, "Update", 1)
	})
}

func TestHandler_CounterHandler(t *testing.T) {
	t.Run("1. Should update metric", func(t *testing.T) {
		m := mocks.CounterService{Mock: mock.Mock{}}
		m.On("Increase", "PollCount", metric.Counter(25)).Return(nil)

		h := Handler{Counter: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/counter/PollCount/25", nil)
		h.Post("/update/counter/{name}/{value}", h.CounterHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusOK)
		m.AssertNumberOfCalls(t, "Increase", 1)
	})

	t.Run("2. Should return error when value is float", func(t *testing.T) {
		m := mocks.CounterService{Mock: mock.Mock{}}
		m.On("Increase", mock.Anything, mock.Anything).Return(nil)

		h := Handler{Counter: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/counter/PollCount/0.123", nil)
		h.Post("/update/counter/{name}/{value}", h.CounterHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusBadRequest)
		m.AssertNumberOfCalls(t, "Increase", 0)
	})

	t.Run("3. Should return error when add value return error", func(t *testing.T) {
		m := mocks.CounterService{Mock: mock.Mock{}}
		m.On("Increase", mock.Anything, mock.Anything).Return(errors.New("error"))

		h := Handler{Counter: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/counter/PollCount/25", nil)

		w := httptest.NewRecorder()
		h.Post("/update/counter/{name}/{value}", h.CounterHandler())

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusBadRequest)
		m.AssertNumberOfCalls(t, "Increase", 1)
	})
}
