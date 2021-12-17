package metrics

import "runtime"

type nextGC struct {
	gaugeBaseMetric
}

func (a nextGC) Name() string {
	return "NextGC"
}

func (a nextGC) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.NextGC)
}
