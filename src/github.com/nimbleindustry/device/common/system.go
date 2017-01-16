package common

import (
	"syscall"
	"time"

	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

// SystemState defines consumed system resources in time
type SystemState struct {
	Timestamp      time.Time `json:"timeStamp"`
	MemoryConsumed float64   `json:"memoryConsumed"`
	DiskConsumed   float64   `json:"diskConsumed"`
	LoadAverage    float64   `json:"loadAverage"`
}

func systemDiskConsumed() float64 {
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs("/", &fs); err != nil {
		return 0.0
	}
	free := float64(fs.Bfree)
	total := float64(fs.Blocks)
	return 1.0 - (free / total)
}

func systemMemoryConsumed() float64 {
	v, _ := mem.VirtualMemory()
	return v.UsedPercent / 100
}

func systemLoadAverage() float64 {
	avg, _ := load.Avg()
	return avg.Load1
}

// GetSystemState returns consumed system resources in time
func GetSystemState() (state *SystemState) {
	state = &SystemState{
		time.Now().UTC(),
		systemMemoryConsumed(),
		systemDiskConsumed(),
		systemLoadAverage(),
	}
	return
}
