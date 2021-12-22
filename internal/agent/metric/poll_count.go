package metric

import "github.com/Jokcik/praktikum-go-devops/internal/metric"

type PollCount struct {
	metric.CounterBaseMetric
	value metric.Counter
}

func (a *PollCount) Name() string {
	return "PollCount"
}

func (a *PollCount) GetValue() interface{} {
	a.value += 1

	return a.value
}
