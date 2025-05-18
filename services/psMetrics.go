package services

import (
	"sync"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// MemoryStats holds RAM usage data
type MemoryStats struct {
	Total       uint64
	Available   uint64
	Used        uint64
	UsedPercent float64
}

// CPUStats holds CPU usage data
type CPUStats struct {
	Percent float64
}

// Temperature represents temperature data from WMI
type Temperature struct {
	Name        string
	Temperature uint32
}

type MetricsSnapshot struct {
	Timestamp    time.Time
	CPU          CPUStats
	Memory       MemoryStats
	Temperatures []Temperature
}

// Almacenamiento en memoria de snapshots
var (
	metricsHistory []MetricsSnapshot
	historyMutex   sync.Mutex
	maxHistory     = 1000 // puedes ajustar este valor según tus necesidades
)

// Guarda un snapshot actual en el historial
func StoreCurrentMetricsSnapshot() error {
	cpuStats, err := GetCPUUsage()
	if err != nil {
		return err
	}
	memStats, err := GetMemoryUsage()
	if err != nil {
		return err
	}
	temps, _ := GetTemperatures()

	snapshot := MetricsSnapshot{
		Timestamp:    time.Now(),
		CPU:          cpuStats,
		Memory:       memStats,
		Temperatures: temps,
	}

	historyMutex.Lock()
	defer historyMutex.Unlock()
	metricsHistory = append(metricsHistory, snapshot)
	if len(metricsHistory) > maxHistory {
		metricsHistory = metricsHistory[1:]
	}
	return nil
}

// Llama a esta función periódicamente (por ejemplo, cada minuto) desde un goroutine en main.go
func StartMetricsCollector(interval time.Duration, stopCh <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			StoreCurrentMetricsSnapshot()
		case <-stopCh:
			return
		}
	}
}

// Devuelve los snapshots más cercanos a hace 10, 5 minutos y el actual
func GetHistoricalMetrics() ([]MetricsSnapshot, error) {
	historyMutex.Lock()
	defer historyMutex.Unlock()

	now := time.Now()
	targets := []time.Duration{10 * time.Minute, 5 * time.Minute, 0}
	result := make([]MetricsSnapshot, 0, 3)

	for _, target := range targets {
		var closest *MetricsSnapshot
		minDiff := time.Duration(1<<63 - 1)
		for i := range metricsHistory {
			diff := absDuration(metricsHistory[i].Timestamp.Sub(now.Add(-target)))
			if diff < minDiff {
				minDiff = diff
				closest = &metricsHistory[i]
			}
		}
		if closest != nil {
			result = append(result, *closest)
		}
	}
	return result, nil
}

func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

// GetCPUUsage returns the CPU usage percent over an interval
func GetCPUUsage() (CPUStats, error) {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return CPUStats{}, err
	}
	return CPUStats{Percent: percentages[0]}, nil
}

// GetMemoryUsage returns memory usage stats
func GetMemoryUsage() (MemoryStats, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return MemoryStats{}, err
	}
	return MemoryStats{
		Total:       vm.Total,
		Available:   vm.Available,
		Used:        vm.Used,
		UsedPercent: vm.UsedPercent,
	}, nil
}

// GetTemperatures fetches CPU temperatures via WMI
func GetTemperatures() ([]Temperature, error) {
	var dst []struct {
		CurrentTemperature uint32
		InstanceName       string
	}
	query := "SELECT CurrentTemperature, InstanceName FROM MSAcpi_ThermalZoneTemperature"
	err := wmi.Query(query, &dst)
	if err != nil {
		return nil, err
	}

	temps := make([]Temperature, len(dst))
	for i, v := range dst {
		// WMI temp is in tenths of Kelvin
		kelvin := float64(v.CurrentTemperature) / 10
		celsius := uint32(kelvin - 273.15)
		temps[i] = Temperature{
			Name:        v.InstanceName,
			Temperature: celsius,
		}
	}
	return temps, nil
}
