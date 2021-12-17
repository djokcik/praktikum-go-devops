package metrics

import "runtime"

type lookups struct {
	gaugeBaseMetric
}

func (a lookups) Name() string {
	return "Lookups"
}

func (a lookups) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.Lookups)
}
