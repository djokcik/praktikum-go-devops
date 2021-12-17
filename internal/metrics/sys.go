package metrics

import "runtime"

type sys struct {
	gaugeBaseMetric
}

func (a sys) Name() string {
	return "Sys"
}

func (a sys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.Sys)
}
