package storage

import "fmt"

type MetricRepository struct {
	BaseRepository
}

func (r *MetricRepository) Update(id interface{}, entity interface{}) (bool, error) {
	r.db.Table[id.(string)] = entity

	return true, nil
}

type MetricElement struct {
	Name  string
	Value interface{}
}

func (r *MetricRepository) List() (interface{}, error) {
	var metricList []MetricElement

	for metricType, metricValue := range r.db.Table {
		metricList = append(metricList, MetricElement{Name: metricType, Value: metricValue})
	}

	return metricList, nil
}

func (r *MetricRepository) Get(metricType string) (interface{}, error) {
	value, ok := r.db.Table[metricType]

	if !ok {
		return nil, fmt.Errorf("the metric `%v` didn`t find", metricType)
	}

	return value, nil
}
