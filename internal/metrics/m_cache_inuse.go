package metrics

import "runtime"

type mCacheInuse struct {
	gaugeBaseMetric
}

func (a mCacheInuse) Name() string {
	return "MCacheInuse"
}

func (a mCacheInuse) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.MCacheInuse)
}
