package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/shirou/gopsutil/mem"
)

type FreeMemory struct {
	metric.GaugeBaseMetric
}

func (a *FreeMemory) Name() string {
	return "FreeMemory"
}

func (a *FreeMemory) GetValue() interface{} {
	v, _ := mem.VirtualMemory()

	return metric.Gauge(v.Free)
}
