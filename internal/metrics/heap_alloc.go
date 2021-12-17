package metrics

import "runtime"

type heapAlloc struct {
	gaugeBaseMetric
}

func (a heapAlloc) Name() string {
	return "HeapAlloc"
}

func (a heapAlloc) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.HeapAlloc)
}
