package agent

import "context"

func (a *agent) CollectMetrics(ctx context.Context) func() {
	return func() {
		for _, metric := range a.metrics {
			if ctx.Err() != nil {
				return
			}

			name := metric.Name()

			a.CollectedMetric[name] = SendAgentMetric{Name: name, Type: metric.Type(), Value: metric.GetValue()}
		}
	}
}
