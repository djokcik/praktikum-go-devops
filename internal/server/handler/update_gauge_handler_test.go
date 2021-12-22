package handler

import (
	"errors"
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGaugeHandler(t *testing.T) {
	t.Run("1. Should update metric", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", "Alloc", metric.Gauge(0.123)).Return(true, nil)
		m.On("Get", storage.GetRepositoryFilter{
			Name:         "Alloc",
			Type:         metric.GaugeType,
			DefaultValue: metric.Gauge(0),
		}).Return(metric.Gauge(0), nil)

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/0.123", nil)
		h.Post("/update/gauge/{name}/{value}", h.GaugeHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusOK)
		m.AssertNumberOfCalls(t, "Get", 1)
		m.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("2. Should return error when value is string", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", mock.Anything, mock.Anything).Return(true, nil)
		m.On("Get", storage.GetRepositoryFilter{
			Name:         "Alloc",
			Type:         metric.GaugeType,
			DefaultValue: metric.Gauge(0),
		}).Return(metric.Gauge(0), nil)

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/test", nil)
		h.Post("/update/gauge/{name}/{value}", h.GaugeHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusBadRequest)
		m.AssertNumberOfCalls(t, "Get", 0)
		m.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("3. Should return error when update was error", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", mock.Anything, mock.Anything).Return(false, errors.New("error"))
		m.On("Get", storage.GetRepositoryFilter{
			Name:         "Alloc",
			Type:         metric.GaugeType,
			DefaultValue: metric.Gauge(0),
		}).Return(metric.Gauge(0), nil)

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/0.123", nil)
		h.Post("/update/gauge/{name}/{value}", h.GaugeHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusBadRequest)
		m.AssertNumberOfCalls(t, "Get", 1)
		m.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("4. Should add value to old metric", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", "Alloc", metric.Gauge(0.444)).Return(true, nil)
		m.On("Get", storage.GetRepositoryFilter{
			Name:         "Alloc",
			Type:         metric.GaugeType,
			DefaultValue: metric.Gauge(0),
		}).Return(metric.Gauge(0.321), nil)

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/0.123", nil)
		h.Post("/update/gauge/{name}/{value}", h.GaugeHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusOK)
		m.AssertNumberOfCalls(t, "Get", 1)
		m.AssertNumberOfCalls(t, "Update", 1)
	})
}
