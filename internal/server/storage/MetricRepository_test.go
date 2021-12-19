package storage

import (
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetricRepository_Update(t *testing.T) {
	t.Run("Should update repository table", func(t *testing.T) {
		db := new(model.Database)

		repository := new(MetricRepository)
		repository.Configure(db)

		_, err := repository.Update("MetricName", 0.123)
		if err != nil {
			t.Errorf("error update repository: %v", err)
		}

		require.Equal(t, db.Table["MetricName"], 0.123)
	})
}
