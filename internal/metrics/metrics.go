package metrics

//go:generate mockery --name=Metric

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
		new(Alloc),
		new(BuckHashSys),
		new(Frees),
		new(GCSys),
		new(GCCPUFraction),
		new(HeapAlloc),
		new(HeapIdle),
		new(HeapInuse),
		new(HeapObjects),
		new(HeapReleased),
		new(HeapSys),
		new(LastGC),
		new(Lookups),
		new(MCacheInuse),
		new(MCacheSys),
		new(MSpanInuse),
		new(MSpanSys),
		new(NextGC),
		new(NumForcedGC),
		new(NumGC),
		new(OtherSys),
		new(PauseTotalNs),
		new(StackInuse),
		new(StackSys),
		new(Sys),
		new(RandomValue),
	}
}

func GetCounterMetrics() []Metric {
	return []Metric{
		new(PollCount),
	}
}
