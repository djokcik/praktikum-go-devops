package metrics

import "runtime"

type lastGC struct {
	gaugeBaseMetric
}

func (a lastGC) Name() string {
	return "LastGC"
}

func (a lastGC) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.LastGC)
}
