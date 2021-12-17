package metrics

import "runtime"

type heapSys struct {
	gaugeBaseMetric
}

func (a heapSys) Name() string {
	return "HeapSys"
}

func (a heapSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.HeapSys)
}
