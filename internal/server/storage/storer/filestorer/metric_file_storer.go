package filestorer

import (
	"context"
	"encoding/json"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/storer"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
	"os"
)

type inMemoryMetricDB struct {
	CounterMapMetric map[string]metric.Counter
	GaugeMapMetric   map[string]metric.Gauge
}

type storeEvent struct {
	CounterMapMetric map[string]metric.Counter
	GaugeMapMetric   map[string]metric.Gauge
}

type metricFileStoreReader struct {
	file    *os.File
	decoder *json.Decoder
}

func newMetricFileStoreReader(filename string) (*metricFileStoreReader, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}

	return &metricFileStoreReader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (r *metricFileStoreReader) ReadEvent() (*storeEvent, error) {
	var event storeEvent
	err := r.decoder.Decode(&event)

	return &event, err
}

func (r *metricFileStoreReader) Close() error {
	return r.file.Close()
}

type metricFileStoreWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func newMetricFileStoreWriter(filename string) (*metricFileStoreWriter, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}

	return &metricFileStoreWriter{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (w *metricFileStoreWriter) SaveEvent(event *storeEvent) error {
	return w.encoder.Encode(&event)
}

func (w *metricFileStoreWriter) Close() error {
	return w.file.Close()
}

type MetricFileStorer struct {
	inMemoryDB *inMemoryMetricDB
	cfg        server.Config

	FileReader *metricFileStoreReader
	FileWriter *metricFileStoreWriter
}

func NewMetricFileStorer(ctx context.Context, cfg server.Config) storer.MetricStorer {
	s := &MetricFileStorer{inMemoryDB: &inMemoryMetricDB{}, cfg: cfg}
	filename := cfg.StoreFile

	reader, err := newMetricFileStoreReader(filename)
	if err != nil {
		s.Log(ctx).Fatal().Err(err).Msg("failed to open reader")
	}

	writer, err := newMetricFileStoreWriter(filename)
	if err != nil {
		s.Log(ctx).Fatal().Err(err).Msg("failed to open writer")
	}

	s.FileReader = reader
	s.FileWriter = writer

	return s
}

func (s *MetricFileStorer) SaveDBValue(ctx context.Context) {
	s.Log(ctx).Info().Msg("save metrics to file")

	err := s.FileWriter.file.Truncate(0)
	if err != nil {
		s.Log(ctx).Error().Err(err).Msgf("invalid truncate metrics")
		return
	}

	_, err = s.FileWriter.file.Seek(0, 0)
	if err != nil {
		s.Log(ctx).Error().Err(err).Msgf("invalid update seek")
		return
	}

	err = s.FileWriter.encoder.Encode(&storeEvent{GaugeMapMetric: s.inMemoryDB.GaugeMapMetric, CounterMapMetric: s.inMemoryDB.CounterMapMetric})
	if err != nil {
		s.Log(ctx).Error().Err(err).Msgf("invalid save metrics")
		return
	}
}

func (s *MetricFileStorer) RestoreDBValue(ctx context.Context) {
	event, err := s.FileReader.ReadEvent()
	if err == nil {
		s.inMemoryDB.GaugeMapMetric = event.GaugeMapMetric
		s.inMemoryDB.CounterMapMetric = event.CounterMapMetric

		s.Log(ctx).Info().Msgf("metrics restored from file %s", s.cfg.StoreFile)
	}
}

func (s *MetricFileStorer) SetCounterDB(counterMapMetric map[string]metric.Counter) {
	s.inMemoryDB.CounterMapMetric = counterMapMetric
}

func (s *MetricFileStorer) SetGaugeDB(gaugeMapMetric map[string]metric.Gauge) {
	s.inMemoryDB.GaugeMapMetric = gaugeMapMetric
}

func (s *MetricFileStorer) Close() {
	s.FileReader.Close()
	s.FileWriter.Close()
}

func (s *MetricFileStorer) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "MetricFileStorer").Logger()

	return &logger
}
