package services

import (
	//"fmt"
	"time"

	"github.com/StackExchange/wmi"
	//"github.com/mindprince/gonvml"
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

/* GPUStats holds GPU usage and temperature
type GPUStats struct {
	Name        string
	Utilization uint
	Temperature uint
}*/

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

/*
// GetGPUStats uses NVML to fetch GPU utilization and temperature
func GetGPUStats() ([]GPUStats, error) {
	if !gonvml.IsCgoEnabled() {
		return nil, fmt.Errorf("NVML is disabled: binary built without CGO")
	}

	if err := gonvml.Initialize(); err != nil {
		return nil, err
	}
	defer gonvml.Shutdown()

	count, err := gonvml.DeviceCount()
	if err != nil {
		return nil, err
	}

	var gpus []GPUStats
	for i := uint(0); i < count; i++ {
		dev, err := gonvml.DeviceHandleByIndex(i)
		if err != nil {
			continue
		}
		name, _ := dev.Name()
		util, _, _ := dev.UtilizationRates()
		temp, _ := dev.Temperature()
		gpus = append(gpus, GPUStats{
			Name:        name,
			Utilization: util,
			Temperature: temp,
		})
	}
	return gpus, nil
}*/
