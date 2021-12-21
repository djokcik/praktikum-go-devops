package handler

import (
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

func TestHandler_ListHandler(t *testing.T) {
	t.Run("1. Should return list metrics in html", func(t *testing.T) {
		metricList := []storage.MetricElement{{Name: "TestType", Value: "TestValue"}}

		m := mocks.Repository{Mock: mock.Mock{}}
		m.On("List").Return(metricList, nil)

		h := Handler{Repo: &m, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		h.Get("/", h.ListHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		m.AssertNumberOfCalls(t, "List", 1)
		require.Equal(t, res.StatusCode, http.StatusOK)
		require.Equal(t, string(resBody), "<html><body><ul><li>Name: TestType<br>Value: TestValue</li></ul></body></html>")
	})
}
