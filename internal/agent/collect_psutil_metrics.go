package agent

import (
	"context"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/google/uuid"
)

func (a *agent) CollectPsutilMetrics(ctx context.Context) func() {
	return func() {
		traceID, _ := uuid.NewUUID()
		logger := a.Log(ctx).With().Str(logging.ServiceKey, "CollectPsutilMetrics").Str(logging.TraceIDKey, traceID.String()).Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		a.Log(ctx).Info().Msg("CollectPsutil")
		a.Lock()
		for _, metric := range GetAgentPsutilMetrics() {
			if ctx.Err() != nil {
				a.Log(ctx).Err(ctx.Err())
				return
			}

			name := metric.Name()

			a.CollectedMetric[name] = SendAgentMetric{Name: name, Type: metric.Type(), Value: metric.GetValue()}
		}
		a.Unlock()
		a.Log(ctx).Info().Msg("Finish collect psutil")
	}
}
