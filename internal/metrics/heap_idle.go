package metrics

import "runtime"

type heapIdle struct {
	gaugeBaseMetric
}

func (a heapIdle) Name() string {
	return "HeapIdle"
}

func (a heapIdle) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.HeapIdle)
}
