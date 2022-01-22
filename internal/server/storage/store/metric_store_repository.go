package store

import (
	"context"
	"encoding/json"
	"github.com/djokcik/praktikum-go-devops/internal/helpers"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/model"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
	"os"
	"sync"
)

//go:generate mockery --name=MetricStore

type MetricStore interface {
	Configure(ctx context.Context, wg *sync.WaitGroup)
	NotifyUpdateDBValue(ctx context.Context)
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

type MetricStoreFile struct {
	DB  *model.Database
	Cfg server.Config

	FileReader *metricFileStoreReader
	FileWriter *metricFileStoreWriter
}

func (s *MetricStoreFile) Configure(ctx context.Context, wg *sync.WaitGroup) {
	filename := s.Cfg.StoreFile
	if filename == "" {
		s.Log(ctx).Info().Msg("save metrics to file are disabled")
		return
	}

	wg.Add(1)
	go func() {
		<-ctx.Done()
		defer wg.Done()

		s.SaveDBValue(ctx)

		s.FileReader.Close()
		s.FileWriter.Close()
	}()

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

	if s.Cfg.Restore {
		s.RestoreDBValue(ctx)
	}

	if s.Cfg.StoreInterval != 0 {
		go helpers.SetTicker(func() { s.SaveDBValue(ctx) }, s.Cfg.StoreInterval)
	}
}

func (s *MetricStoreFile) NotifyUpdateDBValue(ctx context.Context) {
	if s.Cfg.StoreInterval == 0 {
		s.SaveDBValue(ctx)
	}
}

func (s *MetricStoreFile) SaveDBValue(ctx context.Context) {
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

	err = s.FileWriter.encoder.Encode(&storeEvent{GaugeMapMetric: s.DB.GaugeMapMetric, CounterMapMetric: s.DB.CounterMapMetric})
	if err != nil {
		s.Log(ctx).Error().Err(err).Msgf("invalid save metrics")
		return
	}
}

func (s *MetricStoreFile) RestoreDBValue(ctx context.Context) {
	event, err := s.FileReader.ReadEvent()
	if err == nil {
		s.DB.GaugeMapMetric = event.GaugeMapMetric
		s.DB.CounterMapMetric = event.CounterMapMetric

		s.Log(ctx).Info().Msgf("metrics restored from file %s", s.Cfg.StoreFile)
	}
}

func (s *MetricStoreFile) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "MetricStoreFile").Logger()

	return &logger
}
