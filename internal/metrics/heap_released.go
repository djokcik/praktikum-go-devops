package metrics

import "runtime"

type heapReleased struct {
	gaugeBaseMetric
}

func (a heapReleased) Name() string {
	return "HeapReleased"
}

func (a heapReleased) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.HeapReleased)
}
