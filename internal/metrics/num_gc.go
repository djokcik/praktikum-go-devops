package metrics

import "runtime"

type numGC struct {
	gaugeBaseMetric
}

func (a numGC) Name() string {
	return "NumGC"
}

func (a numGC) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.NumGC)
}
