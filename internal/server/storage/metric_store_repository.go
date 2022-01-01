package storage

import (
	"encoding/json"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/helpers"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/model"
	"log"
	"os"
)

type MetricStore interface {
	Configure()
	NotifyUpdateDBValue()
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

// Close TODO: нужно прокинуть контекст и закрыть файл
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

// Close TODO: нужно прокинуть контекст и закрыть файл
func (w *metricFileStoreWriter) Close() error {
	return w.file.Close()
}

type MetricStoreFile struct {
	DB  *model.Database
	Cfg *server.Config

	FileReader *metricFileStoreReader
	FileWriter *metricFileStoreWriter
}

func (s *MetricStoreFile) Configure() {
	filename := s.Cfg.StoreFile
	if filename == "" {
		log.Printf("Save metrics to file are disabled")
		return
	}

	reader, err := newMetricFileStoreReader(filename)
	if err != nil {
		log.Fatal(err)
	}

	writer, err := newMetricFileStoreWriter(filename)
	if err != nil {
		log.Fatal(err)
	}

	s.FileReader = reader
	s.FileWriter = writer

	if s.Cfg.Restore {
		s.RestoreDBValue()
	}

	if s.Cfg.StoreInterval != 0 {
		go helpers.SetTicker(s.SaveDBValue, s.Cfg.StoreInterval)
	}
}

func (s *MetricStoreFile) NotifyUpdateDBValue() {
	if s.Cfg.StoreInterval == 0 {
		s.SaveDBValue()
	}
}

func (s *MetricStoreFile) SaveDBValue() {
	err := s.FileWriter.file.Truncate(0)
	if err != nil {
		log.Printf("invalid truncate metrics with error: %s", err)
		return
	}

	_, err = s.FileWriter.file.Seek(0, 0)
	if err != nil {
		log.Printf("invalid update seek with error: %s", err)
		return
	}

	err = s.FileWriter.encoder.Encode(&storeEvent{GaugeMapMetric: s.DB.GaugeMapMetric, CounterMapMetric: s.DB.CounterMapMetric})
	if err != nil {
		log.Printf("invalid save metrics with error: %s", err)
		return
	}
}

func (s *MetricStoreFile) RestoreDBValue() {
	event, err := s.FileReader.ReadEvent()
	if err == nil {
		s.DB.GaugeMapMetric = event.GaugeMapMetric
		s.DB.CounterMapMetric = event.CounterMapMetric

		log.Println(fmt.Sprintf("metrics restored from file %s", s.Cfg.StoreFile))
	}
}
