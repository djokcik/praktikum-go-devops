package metrics

type Gauge float64
type Counter int64

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type GaugeBaseMetric struct {
}

func (v GaugeBaseMetric) Type() string {
	return GaugeType
}

type CounterBaseMetric struct {
}

func (v CounterBaseMetric) Type() string {
	return CounterType
}

type Metric interface {
	Type() string
	Name() string
}

func GetGaugesMetrics() []Metric {
	return []Metric{
		BuckHashSys{},
		Frees{},
		GCSys{},
		GCCPUFraction{},
		HeapAlloc{},
		HeapIdle{},
		HeapInuse{},
		HeapObjects{},
		HeapReleased{},
		HeapSys{},
		LastGC{},
		Lookups{},
		MCacheInuse{},
		MCacheSys{},
		MSpanInuse{},
		MSpanSys{},
		NextGC{},
		NumForcedGC{},
		NumGC{},
		OtherSys{},
		PauseTotalNs{},
		StackInuse{},
		StackSys{},
		Sys{},
	}
}

func GetCounterMetrics() []Metric {
	return []Metric{
		PollCount{},
	}
}
