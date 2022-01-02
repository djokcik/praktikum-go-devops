package agent

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/google/uuid"
)

func (a *agent) CollectMetrics(ctx context.Context) func() {
	return func() {
		traceID, _ := uuid.NewUUID()
		logger := a.Log(ctx).With().Str(logging.ServiceKey, "CollectMetrics").Str(logging.TraceIDKey, traceID.String()).Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		a.Log(ctx).Info().Msg("Collect")
		for _, metric := range a.metrics {
			if ctx.Err() != nil {
				a.Log(ctx).Err(ctx.Err())
				return
			}

			name := metric.Name()

			a.CollectedMetric[name] = SendAgentMetric{Name: name, Type: metric.Type(), Value: metric.GetValue()}
		}
		a.Log(ctx).Info().Msg("Finish collect")
	}
}
