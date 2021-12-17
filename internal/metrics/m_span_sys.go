package metrics

import "runtime"

type mSpanSys struct {
	gaugeBaseMetric
}

func (a mSpanSys) Name() string {
	return "MSpanSys"
}

func (a mSpanSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.MSpanSys)
}
