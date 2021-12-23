package agent

func (a *agent) CollectMetrics() {
	for _, metric := range a.metrics {
		if a.Ctx.Err() != nil {
			return
		}

		name := metric.Name()

		a.CollectedMetric[name] = SendAgentMetric{Name: name, Type: metric.Type(), Value: metric.GetValue()}
	}
}
