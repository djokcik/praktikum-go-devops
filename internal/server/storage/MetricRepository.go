package storage

type MetricRepository struct {
	BaseRepository
}

func (r *MetricRepository) Update(id interface{}, entity interface{}) (bool, error) {
	r.db.Table[id.(string)] = entity

	return true, nil
}
