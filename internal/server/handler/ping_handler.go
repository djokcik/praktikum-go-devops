package handler

import (
	"github.com/djokcik/praktikum-go-devops/internal/server/storage"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"net/http"
)

func (h *Handler) PingHandler(repository *storage.MetricDatabaseRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "Ping").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		err := repository.Ping()
		if err != nil {
			h.Log(ctx).Error().Err(err).Msg("failed connect to database")
			http.Error(rw, "failed to connect", http.StatusInternalServerError)
			return
		}

		rw.Write([]byte("OK"))
	}
}
