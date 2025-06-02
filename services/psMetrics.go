package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type Metric struct {
	Timestamp   time.Time
	CPUUsage    float64 `json:"cpu_usage"`
	RAMUsage    float64 `json:"ram_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	Temperature float64 `json:"temperature"`
}

var (
	metrics []Metric
	mu      sync.Mutex
)

func collectMetrics(interval time.Duration) {
	for {
		cpuPercent, _ := cpu.Percent(0, false)
		vmStat, _ := mem.VirtualMemory()
		diskStat, _ := disk.Usage("C:\\\\") // "/" en Linux, "C:\\" en Windows

		temps, _ := host.SensorsTemperatures()
		var temp float64
		if len(temps) > 0 {
			for _, t := range temps {
				if t.Temperature > 0 {
					temp = t.Temperature
					break
				}
			}
		}

		metric := Metric{
			Timestamp:   time.Now(),
			CPUUsage:    cpuPercent[0],
			RAMUsage:    vmStat.UsedPercent,
			DiskUsage:   diskStat.UsedPercent,
			Temperature: temp,
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

func averageUsage(duration time.Duration) (cpuAvg, ramAvg, diskAvg, tempAvg float64) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	count := 0
	for _, m := range metrics {
		if now.Sub(m.Timestamp) <= duration {
			cpuAvg += m.CPUUsage
			ramAvg += m.RAMUsage
			diskAvg += m.DiskUsage
			tempAvg += m.Temperature
			count++
		}
	}

	if count > 0 {
		cpuAvg /= float64(count)
		ramAvg /= float64(count)
		diskAvg /= float64(count)
		tempAvg /= float64(count)
	}

	return
}

func printMetrics() {
	for {
		cpu5, ram5, disk5, temp5 := averageUsage(5 * time.Minute)
		cpu10, ram10, disk10, temp10 := averageUsage(10 * time.Minute)

		mu.Lock()
		var current Metric
		if len(metrics) > 0 {
			current = metrics[len(metrics)-1]
		}
		mu.Unlock()

		fmt.Printf("\n=== MÉTRICAS ===\n")
		fmt.Printf("Actual   → CPU: %.2f%% | RAM: %.2f%% | Disco: %.2f%% | Temp: %.2f°C\n",
			current.CPUUsage, current.RAMUsage, current.DiskUsage, current.Temperature)
		fmt.Printf("Hace 5m  → CPU: %.2f%% | RAM: %.2f%% | Disco: %.2f%% | Temp: %.2f°C\n",
			cpu5, ram5, disk5, temp5)
		fmt.Printf("Hace 10m → CPU: %.2f%% | RAM: %.2f%% | Disco: %.2f%% | Temp: %.2f°C\n",
			cpu10, ram10, disk10, temp10)

		time.Sleep(10 * time.Second)
	}
}
