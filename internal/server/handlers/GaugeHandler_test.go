package handlers

import (
	"errors"
	"github.com/Jokcik/praktikum-go-devops/internal/metrics"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGaugeHandler(t *testing.T) {
	t.Run("1. Should update metric", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", "Alloc", metrics.Gauge(0.123)).Return(true, nil)

		request := httptest.NewRequest(http.MethodPost, "/update/counter/Alloc/0.123", nil)

		w := httptest.NewRecorder()
		h := GaugeHandler(&m)

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusOK)
		m.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("2. Should return error when value is string", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", mock.Anything, mock.Anything).Return(true, nil)

		request := httptest.NewRequest(http.MethodPost, "/update/counter/Alloc/test", nil)

		w := httptest.NewRecorder()
		h := GaugeHandler(&m)

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusBadRequest)
		m.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("3. Should return error when update was error", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", mock.Anything, mock.Anything).Return(false, errors.New("error"))

		request := httptest.NewRequest(http.MethodPost, "/update/counter/Alloc/0.123", nil)

		w := httptest.NewRecorder()
		h := GaugeHandler(&m)

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusBadRequest)
		m.AssertNumberOfCalls(t, "Update", 1)
	})
}
