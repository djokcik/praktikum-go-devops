package metrics

import "runtime"

type buckHashSys struct {
	gaugeBaseMetric
}

func (a buckHashSys) Name() string {
	return "BuckHashSys"
}

func (a buckHashSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return gauge(memStats.BuckHashSys)
}
