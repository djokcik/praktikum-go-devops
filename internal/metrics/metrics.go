package metrics

type gauge float64
type counter int64

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type gaugeBaseMetric struct {
}

func (v gaugeBaseMetric) Type() string {
	return Gauge
}

type counterBaseMetric struct {
}

func (v counterBaseMetric) Type() string {
	return Counter
}

type Metric interface {
	Type() string
	Name() string
	GetValue() interface{}
}

func GetMetrics() []Metric {
	return []Metric{
		// gauges
		buckHashSys{},
		frees{},
		gCSys{},
		gCCPUFraction{},
		heapAlloc{},
		heapIdle{},
		heapInuse{},
		heapObjects{},
		heapReleased{},
		heapSys{},
		lastGC{},
		lookups{},
		mCacheInuse{},
		mCacheSys{},
		mSpanInuse{},
		mSpanSys{},
		nextGC{},
		numForcedGC{},
		numGC{},
		otherSys{},
		pauseTotalNs{},
		stackInuse{},
		stackSys{},
		sys{},

		// counter
		PollCount{},
	}
}
