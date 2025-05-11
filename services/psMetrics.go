package services

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type Metric struct {
	TimeStamp   time.Time `json:"timestamp"`
	CPUUsage    float64   `json:"cpuusage"`
	RAMUsage    float64   `json:"ramusage"`
	GPUUsage    float64   `json:"gpuusage"`
	Temperature float64   `json:"temperature"`
}

var metricsHistory []Metric

func CollectMetrics() {
	for {
		cpuPercent, _ := cpu.Percent(0, false)
		memStats, _ := mem.VirtualMemory()
		temperatures, _ := host.SensorsTemperatures()

		temp := 0.0
		if len(temperatures) > 0 {
			temp = float64(temperatures[0].Temperature)
		}

		m := Metric{
			TimeStamp:   time.Now(),
			CPUUsage:    cpuPercent[0],
			RAMUsage:    memStats.UsedPercent,
			GPUUsage:    0.0,
			Temperature: temp,
		}
		metricsHistory = append(metricsHistory, m)

		if len(metricsHistory) > 600 {
			metricsHistory = metricsHistory[len(metricsHistory)-600:]
		}

		time.Sleep(10 * time.Second)
	}
}

func GetMetrics() map[string]*Metric {
	now := time.Now()
	var m10, m5, m1 *Metric

	for i := len(metricsHistory) - 1; i >= 0; i-- {
		m := metricsHistory[i]
		if m10 == nil && now.Sub(m.TimeStamp) > 10*time.Minute {
			m10 = &m
		}

		if m5 == nil && now.Sub(m.TimeStamp) > 5*time.Minute {
			m5 = &m
		}

		if m1 == nil && now.Sub(m.TimeStamp) > 1*time.Minute {
			m1 = &m
		}
	}

	return map[string]*Metric{
		"10 Minutes Ago: ": m10,
		"5 Minutes Ago: ": m5,
		"1 Minute Ago: ": m1,
		"Current: ": &metricsHistory[len(metricsHistory)-1]
	}
}
