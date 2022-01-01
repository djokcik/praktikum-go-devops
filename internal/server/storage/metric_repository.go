package storage

import (
	"errors"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/model"
)

var (
	ErrValueNotFound = errors.New("value not found")
)

type MetricRepository struct {
	BaseRepository
	Store MetricStore
}

func (r *MetricRepository) Configure(db *model.Database, cfg *server.Config) {
	r.BaseRepository.Configure(db, cfg)

	r.Store = &MetricStoreFile{DB: db, Cfg: cfg}
	r.Store.Configure()
}

func (r *MetricRepository) Update(name interface{}, value interface{}) (bool, error) {
	r.db.Lock()
	defer r.db.Unlock()

	// TODO use db that code will be more simply
	switch metricValue := value.(type) {
	default:
		return false, fmt.Errorf("entity could`t convert `%v` into available metric type", value)
	case metric.Counter:
		r.db.CounterMapMetric[name.(string)] = metricValue
	case metric.Gauge:
		r.db.GaugeMapMetric[name.(string)] = metricValue
	}

	r.Store.NotifyUpdateDBValue()

	return true, nil
}

func (r *MetricRepository) List(filter *ListRepositoryFilter) (interface{}, error) {
	var metricList []metric.Metric

	// TODO use db that code will be more simply
	switch filter.Type {
	default:
		return nil, fmt.Errorf("type `%v` isn`t avalilable metric type", filter.Type)
	case metric.GaugeType:
		for metricName, metricValue := range r.db.GaugeMapMetric {
			metricList = append(metricList, metric.Metric{Name: metricName, Value: metricValue})
		}
	case metric.CounterType:
		for metricName, metricValue := range r.db.CounterMapMetric {
			metricList = append(metricList, metric.Metric{Name: metricName, Value: metricValue})
		}
	}

	return metricList, nil
}

func (r *MetricRepository) Get(filter *GetRepositoryFilter) (interface{}, error) {
	metricType := filter.Type
	var value interface{}
	var ok bool

	// TODO use db that code will be more simply
	switch metricType {
	default:
		return nil, fmt.Errorf("type `%v` isn`t avalilable metric type", filter.Type)
	case metric.GaugeType:
		value, ok = r.db.GaugeMapMetric[filter.Name]
		if !ok {
			return 0, ErrValueNotFound
		}
	case metric.CounterType:
		value, ok = r.db.CounterMapMetric[filter.Name]
		if !ok {
			return 0, ErrValueNotFound
		}
	}

	return value, nil
}
