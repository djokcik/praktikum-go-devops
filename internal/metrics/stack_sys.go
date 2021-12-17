package metrics

import "runtime"

type stackSys struct {
	gaugeBaseMetric
}

func (a stackSys) Name() string {
	return "StackSys"
}

func (a stackSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.StackSys)
}
