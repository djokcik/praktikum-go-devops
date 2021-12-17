package metrics

import "runtime"

type mSpanInuse struct {
	gaugeBaseMetric
}

func (a mSpanInuse) Name() string {
	return "MSpanInuse"
}

func (a mSpanInuse) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.MSpanInuse)
}
