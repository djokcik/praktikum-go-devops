package metrics

import "runtime"

type alloc struct {
	gaugeBaseMetric
}

func (a alloc) Name() string {
	return "Alloc"
}

func (a alloc) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.Alloc)
}
