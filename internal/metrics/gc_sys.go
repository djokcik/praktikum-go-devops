package metrics

import "runtime"

type gCSys struct {
	gaugeBaseMetric
}

func (a gCSys) Name() string {
	return "GCSys"
}

func (a gCSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.GCSys)
}
