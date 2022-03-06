package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/djokcik/praktikum-go-devops/internal/server/handler"
	"github.com/djokcik/praktikum-go-devops/internal/server/service"
	counterrepomock "github.com/djokcik/praktikum-go-devops/internal/server/storage/counterrepo/mocks"
	gaugerepomock "github.com/djokcik/praktikum-go-devops/internal/server/storage/gaugerepo/mocks"
	commonService "github.com/djokcik/praktikum-go-devops/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
)

func BenchmarkUpdateListHandler(b *testing.B) {
	counterRepoMock := counterrepomock.Repository{Mock: mock.Mock{}}
	counterRepoMock.On("UpdateList", mock.Anything, mock.Anything).Return(nil)

	gaugeRepoMock := gaugerepomock.Repository{Mock: mock.Mock{}}
	gaugeRepoMock.On("UpdateList", mock.Anything, mock.Anything).Return(nil)

	counter := service.CounterServiceImpl{Hash: &commonService.HashServiceImpl{}, Repo: &counterRepoMock}
	gauge := service.GaugeServiceImpl{Hash: &commonService.HashServiceImpl{}, Repo: &gaugeRepoMock}

	h := handler.Handler{Mux: chi.NewMux(), Counter: &counter, Gauge: &gauge}
	h.Post("/updates/", h.UpdateListJSONHandler())

	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		requestBody := `[{"id":"TestCounterMetric","type":"counter","delta":10,"hash":"hashCounter"},{"id":"TestGaugeMetric","type":"gauge","value":1.5,"hash":"hashGauge"}]`
		request := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewBufferString(requestBody))
		b.StartTimer()

		h.ServeHTTP(w, request)
	}
}
