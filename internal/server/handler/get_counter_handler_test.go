package handler

import (
	"errors"
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetCounterHandler(t *testing.T) {
	t.Run("1. Should return counter metric", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Get", "TestName").Return(metric.Counter(123), nil)

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodGet, "/value/counter/TestName", nil)
		h.Get("/value/counter/{name}", h.GetCounterHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		m.AssertNumberOfCalls(t, "Get", 1)
		require.Equal(t, res.StatusCode, http.StatusOK)
		require.Equal(t, string(resBody), "123")
	})

	t.Run("2. Should return 404 when metric didn`t find", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Get", "TestName").Return(nil, errors.New("error"))

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodGet, "/value/counter/TestName", nil)
		h.Get("/value/counter/{name}", h.GetCounterHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		m.AssertNumberOfCalls(t, "Get", 1)
		require.Equal(t, res.StatusCode, http.StatusNotFound)
	})

	t.Run("3. Should return error when metric type is Gauge", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Get", "TestName").Return(metric.Gauge(123), nil)

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodGet, "/value/counter/TestName", nil)
		h.Get("/value/counter/{name}", h.GetCounterHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		m.AssertNumberOfCalls(t, "Get", 1)
		require.Equal(t, res.StatusCode, http.StatusNotFound)
	})
}
