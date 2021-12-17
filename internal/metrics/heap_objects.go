package metrics

import "runtime"

type heapObjects struct {
	gaugeBaseMetric
}

func (a heapObjects) Name() string {
	return "HeapObjects"
}

func (a heapObjects) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.HeapObjects)
}
