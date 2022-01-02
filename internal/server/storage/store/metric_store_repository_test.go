package store

import (
	context "context"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/model"
	"github.com/stretchr/testify/require"
	"os"
	"sync"
	"testing"
)

func TestMetricStoreFile_Configure(t *testing.T) {
	t.Run("1. should  initialize fileReader and fileWriter", func(t *testing.T) {
		store := &MetricStoreFile{Cfg: &server.Config{StoreFile: "mocks/testfile.txt"}}
		store.Configure(context.Background(), &sync.WaitGroup{})
		defer os.Remove("mocks/testfile.txt")

		require.NotNil(t, store.FileWriter)
		require.NotNil(t, store.FileReader)
	})

	t.Run("2. should  not initialize fileReader and fileWriter", func(t *testing.T) {
		store := &MetricStoreFile{Cfg: &server.Config{StoreFile: ""}}
		store.Configure(context.Background(), &sync.WaitGroup{})

		require.Nil(t, store.FileWriter)
		require.Nil(t, store.FileReader)
	})

	t.Run("3. should save data after cancel context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		wg := sync.WaitGroup{}
		counterMap := make(map[string]metric.Counter)
		counterMap["testCounter"] = metric.Counter(123)

		gaugeMap := make(map[string]metric.Gauge)
		gaugeMap["testGauge"] = metric.Gauge(0.123)

		db := model.Database{CounterMapMetric: counterMap, GaugeMapMetric: gaugeMap}

		store := &MetricStoreFile{Cfg: &server.Config{StoreFile: "mocks/testfile.txt"}, DB: &db}
		store.Configure(ctx, &wg)
		defer os.Remove("mocks/testfile.txt")

		_, err := store.FileReader.ReadEvent()
		require.NotNil(t, err)

		cancel()
		wg.Wait()

		store.FileReader.file.Seek(0, 0)

		reader, _ := newMetricFileStoreReader("mocks/testfile.txt")
		event, err := reader.ReadEvent()
		defer reader.Close()

		require.Nil(t, err)
		require.Equal(t, event.CounterMapMetric, counterMap)
		require.Equal(t, event.GaugeMapMetric, gaugeMap)
	})

	t.Run("4. should restore data", func(t *testing.T) {
		counterMap := make(map[string]metric.Counter)
		counterMap["testCounter"] = metric.Counter(123)

		gaugeMap := make(map[string]metric.Gauge)
		gaugeMap["testGauge"] = metric.Gauge(0.123)

		event := storeEvent{GaugeMapMetric: gaugeMap, CounterMapMetric: counterMap}

		writer, _ := newMetricFileStoreWriter("mocks/testfile.txt")
		defer os.Remove("mocks/testfile.txt")

		writer.encoder.Encode(event)
		writer.Close()

		db := model.Database{CounterMapMetric: make(map[string]metric.Counter), GaugeMapMetric: make(map[string]metric.Gauge)}

		store := &MetricStoreFile{Cfg: &server.Config{StoreFile: "mocks/testfile.txt", Restore: true}, DB: &db}
		store.Configure(context.Background(), &sync.WaitGroup{})

		require.Equal(t, store.DB, &model.Database{CounterMapMetric: counterMap, GaugeMapMetric: gaugeMap})
	})
}
