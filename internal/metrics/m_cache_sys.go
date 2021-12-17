package metrics

import "runtime"

type mCacheSys struct {
	gaugeBaseMetric
}

func (a mCacheSys) Name() string {
	return "MCacheSys"
}

func (a mCacheSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.MCacheSys)
}
