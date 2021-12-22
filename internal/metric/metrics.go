package metric

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

type Metric struct {
	Name  string
	Value interface{}
}
