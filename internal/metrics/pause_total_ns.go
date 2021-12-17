package metrics

import "runtime"

type pauseTotalNs struct {
	gaugeBaseMetric
}

func (a pauseTotalNs) Name() string {
	return "PauseTotalNs"
}

func (a pauseTotalNs) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.PauseTotalNs)
}
