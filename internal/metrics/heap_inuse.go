package metrics

import "runtime"

type heapInuse struct {
	gaugeBaseMetric
}

func (a heapInuse) Name() string {
	return "HeapInuse"
}

func (a heapInuse) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.HeapInuse)
}
