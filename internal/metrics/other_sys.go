package metrics

import "runtime"

type otherSys struct {
	gaugeBaseMetric
}

func (a otherSys) Name() string {
	return "OtherSys"
}

func (a otherSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.OtherSys)
}
