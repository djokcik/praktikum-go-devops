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

func TestCounterHandler(t *testing.T) {
	t.Run("1. Should update metric", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", "PollCount", metrics.Counter(25)).Return(true, nil)

		request := httptest.NewRequest(http.MethodPost, "/update/counter/PollCount/25", nil)

		w := httptest.NewRecorder()
		h := CounterHandler(&m)

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusOK)
		m.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("2. Should return error when value is float", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", mock.Anything, mock.Anything).Return(true, nil)

		request := httptest.NewRequest(http.MethodPost, "/update/counter/PollCount/0.123", nil)

		w := httptest.NewRecorder()
		h := CounterHandler(&m)

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusBadRequest)
		m.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("3. Should return error when update was error", func(t *testing.T) {
		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("Update", mock.Anything, mock.Anything).Return(false, errors.New("error"))

		request := httptest.NewRequest(http.MethodPost, "/update/counter/PollCount/25", nil)

		w := httptest.NewRecorder()
		h := CounterHandler(&m)

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusBadRequest)
		m.AssertNumberOfCalls(t, "Update", 1)
	})
}
