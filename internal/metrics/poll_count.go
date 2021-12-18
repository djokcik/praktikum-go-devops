package metrics

type PollCount struct {
	CounterBaseMetric
	value Counter
}

func (a PollCount) Name() string {
	return "PollCount"
}

func (a PollCount) GetValue() interface{} {
	a.value += 1

	return a.value
}
