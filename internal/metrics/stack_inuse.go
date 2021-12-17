package metrics

import "runtime"

type stackInuse struct {
	gaugeBaseMetric
}

func (a stackInuse) Name() string {
	return "StackInuse"
}

func (a stackInuse) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.StackInuse)
}
