package filestorer

import (
	context "context"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMetricStoreFile_Configure(t *testing.T) {
	t.Run("1. should  initialize fileReader and fileWriter", func(t *testing.T) {
		store := NewMetricFileStorer(context.Background(), server.Config{StoreFile: "../mocks/testfile.txt"})
		fileStore := store.(*MetricFileStorer)
		defer os.Remove("../mocks/testfile.txt")

		require.NotNil(t, fileStore.FileWriter)
		require.NotNil(t, fileStore.FileReader)
	})

	t.Run("2. should restore data", func(t *testing.T) {
		counterMap := make(map[string]metric.Counter)
		counterMap["testCounter"] = metric.Counter(123)

		gaugeMap := make(map[string]metric.Gauge)
		gaugeMap["testGauge"] = metric.Gauge(0.123)

		event := storeEvent{GaugeMapMetric: gaugeMap, CounterMapMetric: counterMap}

		writer, _ := newMetricFileStoreWriter("../mocks/testfile.txt")
		defer os.Remove("../mocks/testfile.txt")

		writer.encoder.Encode(event)
		writer.Close()

		db := inMemoryMetricDB{CounterMapMetric: make(map[string]metric.Counter), GaugeMapMetric: make(map[string]metric.Gauge)}

		store := NewMetricFileStorer(context.Background(), server.Config{StoreFile: "../mocks/testfile.txt", Restore: true})
		fileStore := store.(*MetricFileStorer)
		fileStore.inMemoryDB = &db
		fileStore.RestoreDBValue(context.Background())

		require.Equal(t, fileStore.inMemoryDB, &inMemoryMetricDB{CounterMapMetric: counterMap, GaugeMapMetric: gaugeMap})
	})
}
