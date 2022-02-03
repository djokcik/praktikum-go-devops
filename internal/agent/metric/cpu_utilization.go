package metric

import (
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
)

type cpuUtilization struct {
	metric.GaugeBaseMetric
	cpu      int
	cpuValue float64
}

func NewCPUUtilization(cpu int, cpuValue float64) *cpuUtilization {
	return &cpuUtilization{cpu: cpu, cpuValue: cpuValue}
}

func (a *cpuUtilization) Name() string {
	return fmt.Sprintf("CPUutilization%d", a.cpu+1)
}

func (a *cpuUtilization) GetValue() interface{} {
	return metric.Gauge(a.cpuValue)
}
