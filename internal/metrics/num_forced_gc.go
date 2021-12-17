package metrics

import "runtime"

type numForcedGC struct {
	gaugeBaseMetric
}

func (a numForcedGC) Name() string {
	return "NumForcedGC"
}

func (a numForcedGC) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.NumForcedGC)
}
