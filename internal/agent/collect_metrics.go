package agent

func (a *agent) CollectMetrics() {
	for _, metric := range a.metrics {
		name := metric.Name()

		a.CollectedMetric[name] = SendAgentMetric{Name: name, Type: metric.Type(), Value: metric.GetValue()}
	}
}
