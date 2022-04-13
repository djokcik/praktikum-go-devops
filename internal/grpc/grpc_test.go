package grpc

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/internal/server/storage/reporegistry"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestMakeGRPCMetricService(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		s := MakeGRPCMetricService(server.Config{}, reporegistry.NewInMem(context.Background(), &sync.WaitGroup{}, server.Config{}))

		require.NotEqual(t, s, nil)
	})
}
