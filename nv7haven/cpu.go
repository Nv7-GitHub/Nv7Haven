package nv7haven

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type CPUData struct {
	CPUUsage       float64 `json:"cpu_usage"`
	MemoryUsage    float64 `json:"memory_usage"`
	CPUTemperature float64 `json:"cpu_temperature"`
	DiskUsage      float64 `json:"disk_usage"`
}

func (n *Nv7Haven) getCPU(c *fiber.Ctx) error {
	cpuData := CPUData{}

	// CPU Usage
	cpuPercentages, err := cpu.Percent(100*time.Millisecond, false)
	if err != nil {
		return err
	}
	cpuData.CPUUsage = cpuPercentages[0]

	// Memory Usage
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	cpuData.MemoryUsage = virtualMemory.UsedPercent

	// CPU Temperature
	temps, err := host.SensorsTemperatures()
	if err != nil {
		return err
	}
	cpuTemp := 0.0
	for _, temp := range temps {
		if temp.SensorKey == "soc_thermal" { // replace with the correct sensor key
			cpuTemp = temp.Temperature
			break
		}
	}
	cpuData.CPUTemperature = cpuTemp

	// Disk Usage
	diskUsage, err := disk.Usage("/")
	if err != nil {
		return err
	}
	cpuData.DiskUsage = diskUsage.UsedPercent

	return c.JSON(cpuData)
}
