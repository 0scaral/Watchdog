package services

import (
	"strconv"
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

		if metric.CPUUsage > 40 {
			SendAlertMail("High CPU Usage Alert\nCPU usage is above 40%: " + strconv.FormatFloat(metric.CPUUsage, 'f', 2, 64) + "%")
			SendAlertTelegram("High CPU Usage Alert\nCPU usage is above 40%: " + strconv.FormatFloat(metric.CPUUsage, 'f', 2, 64) + "%")
		}

		if metric.RAMUsage > 40 {
			SendAlertMail("High RAM Usage Alert\nRAM usage is above 40%: " + strconv.FormatFloat(metric.RAMUsage, 'f', 2, 64) + "%")
			SendAlertTelegram("High RAM Usage Alert\nRAM usage is above 40%: " + strconv.FormatFloat(metric.RAMUsage, 'f', 2, 64) + "%")
		}

		if metric.DiskUsage > 40 {
			SendAlertMail("High Disk Usage Alert\nDisk usage is above 40%: " + strconv.FormatFloat(metric.DiskUsage, 'f', 2, 64) + "%")
			SendAlertTelegram("High Disk Usage Alert\nDisk usage is above 40%: " + strconv.FormatFloat(metric.DiskUsage, 'f', 2, 64) + "%")
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
