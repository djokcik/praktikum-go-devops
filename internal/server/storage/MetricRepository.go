package storage

import "fmt"

type MetricRepository struct {
	BaseRepository
}

func (r *MetricRepository) Update(id interface{}, entity interface{}) (bool, error) {
	fmt.Println(id, entity)
	r.db.Table[id.(string)] = entity

	return true, nil
}
