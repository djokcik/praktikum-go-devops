package logging

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/pkg"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	contextKeyLogger  = pkg.ContextKey("Logger")
	contextKeyTraceID = pkg.ContextKey("TraceID")
)

func GetCtxLogger(ctx context.Context) (context.Context, zerolog.Logger) {
	if ctxValue := ctx.Value(contextKeyLogger); ctxValue != nil {
		if ctxLogger, ok := ctxValue.(zerolog.Logger); ok {
			return ctx, ctxLogger
		}
	}

	traceID, _ := uuid.NewUUID()
	logger := NewLogger().With().Str(TraceIDKey, traceID.String()).Logger()

	ctx = context.WithValue(ctx, contextKeyTraceID, traceID.String())

	return SetCtxLogger(ctx, logger), logger
}

func SetCtxLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}
