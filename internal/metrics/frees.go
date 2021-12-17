package metrics

import "runtime"

type frees struct {
	gaugeBaseMetric
}

func (a frees) Name() string {
	return "Frees"
}

func (a frees) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.Frees)
}
