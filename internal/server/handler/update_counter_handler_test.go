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

func TestCounterHandler(t *testing.T) {
	t.Run("1. Should update metric", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", "PollCount", metric.Counter(25)).Return(true, nil)
		m.On("Get", storage.GetRepositoryFilter{
			Name:         "PollCount",
			Type:         metric.CounterType,
			DefaultValue: metric.Counter(0),
		}).Return(metric.Counter(0), nil)

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/counter/PollCount/25", nil)
		h.Post("/update/counter/{name}/{value}", h.CounterHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusOK)
		m.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("2. Should return error when value is float", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", mock.Anything, mock.Anything).Return(true, nil)
		m.On("Get", storage.GetRepositoryFilter{
			Name:         "PollCount",
			Type:         metric.CounterType,
			DefaultValue: metric.Counter(0),
		}).Return(metric.Counter(0), nil)

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/counter/PollCount/0.123", nil)
		h.Post("/update/counter/{name}/{value}", h.CounterHandler())

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
			Name:         "PollCount",
			Type:         metric.CounterType,
			DefaultValue: metric.Counter(0),
		}).Return(metric.Counter(0), nil)

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/counter/PollCount/25", nil)

		w := httptest.NewRecorder()
		h.Post("/update/counter/{name}/{value}", h.CounterHandler())

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusBadRequest)
		m.AssertNumberOfCalls(t, "Get", 1)
		m.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("4. Should add value to old metric", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", "PollCount", metric.Counter(125)).Return(true, nil)
		m.On("Get", storage.GetRepositoryFilter{
			Name:         "PollCount",
			Type:         metric.CounterType,
			DefaultValue: metric.Counter(0),
		}).Return(metric.Counter(100), nil)

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodPost, "/update/counter/PollCount/25", nil)
		h.Post("/update/counter/{name}/{value}", h.CounterHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusOK)
		m.AssertNumberOfCalls(t, "Get", 1)
		m.AssertNumberOfCalls(t, "Update", 1)
	})
}
