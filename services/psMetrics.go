package services

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type Metric struct {
	Timestamp time.Time
	CPUUsage  float64 `json:"cpu_usage"`
	RAMUsage  float64 `json:"ram_usage"`
	DiskUsage float64 `json:"disk_usage"`
}

var (
	metrics []Metric
	mu      sync.Mutex
)

func collectMetrics(interval time.Duration) {
	for {
		cpuPercent, _ := cpu.Percent(0, false)
		vmStat, _ := mem.VirtualMemory()
		diskStat, _ := disk.Usage("C:\\\\")

		metric := Metric{
			Timestamp: time.Now(),
			CPUUsage:  cpuPercent[0],
			RAMUsage:  vmStat.UsedPercent,
			DiskUsage: diskStat.UsedPercent,
		}

		mu.Lock()
		metrics = append(metrics, metric)
		cutoff := time.Now().Add(-11 * time.Minute)
		for len(metrics) > 0 && metrics[0].Timestamp.Before(cutoff) {
			metrics = metrics[1:]
		}
		mu.Unlock()

		time.Sleep(interval)
	}
}

func averageUsage(duration time.Duration) (cpuAvg, ramAvg, diskAvg float64) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	count := 0
	for _, m := range metrics {
		if now.Sub(m.Timestamp) <= duration {
			cpuAvg += m.CPUUsage
			ramAvg += m.RAMUsage
			diskAvg += m.DiskUsage
			count++
		}
	}

	if count > 0 {
		cpuAvg /= float64(count)
		ramAvg /= float64(count)
		diskAvg /= float64(count)
	}

	return
}

// GetCurrentMetric retorna la última métrica registrada
func GetCurrentMetric() Metric {
	mu.Lock()
	defer mu.Unlock()
	if len(metrics) > 0 {
		return metrics[len(metrics)-1]
	}
	return Metric{}
}

// GetAverageMetric retorna el promedio de las métricas en el intervalo dado
func GetAverageMetric(duration time.Duration) Metric {
	cpu, ram, disk := averageUsage(duration)
	return Metric{
		CPUUsage:  cpu,
		RAMUsage:  ram,
		DiskUsage: disk,
	}
}

// Inicia la recolección de métricas (llamar desde main)
func StartMetricsCollection() {
	go collectMetrics(10 * time.Second)
}
