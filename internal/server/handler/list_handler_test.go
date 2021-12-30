package handler

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestHandler_ListHandler(t *testing.T) {
	t.Run("1. Should return list metrics in html", func(t *testing.T) {
		// set root repository for load html file
		_, filename, _, _ := runtime.Caller(0)
		dir := path.Join(path.Dir(filename), "../../../") // root repository
		err := os.Chdir(dir)
		if err != nil {
			panic(err)
		}

		metricList := []metric.Metric{{Name: "TestType", Value: "TestValue"}}

		counterMock := mocks.CounterService{Mock: mock.Mock{}}
		counterMock.On("List").Return([]metric.Metric{}, nil)

		gaugeMock := mocks.GaugeService{Mock: mock.Mock{}}
		gaugeMock.On("List").Return(metricList, nil)

		h := Handler{Gauge: &gaugeMock, Counter: &counterMock, Mux: chi.NewMux()}
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		h.Get("/", h.ListHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		counterMock.AssertNumberOfCalls(t, "List", 1)
		gaugeMock.AssertNumberOfCalls(t, "List", 1)

		require.Equal(t, res.StatusCode, http.StatusOK)
		require.NotContains(t, string(resBody), "Counter:")
		require.Contains(t, string(resBody), "Gauges:")
	})
}
