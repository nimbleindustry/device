package common

import (
	"syscall"
	"time"

	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

type SystemState struct {
	Timestamp      time.Time `json:"timeStamp"`
	MemoryConsumed float64   `json:"memoryConsumed"`
	DiskConsumed   float64   `json:"diskConsumed"`
	LoadAverage    float64   `json:"loadAverage"`
}

// SystemDiskConsumed returns the amount (as a percentange) of
// system disk space that has been consumed
func SystemDiskConsumed() float64 {
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs("/", &fs); err != nil {
		return 0.0
	}
	free := float64(fs.Bfree)
	total := float64(fs.Blocks)
	return 1.0 - (free / total)
}

func SystemMemoryConsumed() float64 {
	v, _ := mem.VirtualMemory()
	return v.UsedPercent / 100
}

func SystemLoadAverage() float64 {
	avg, _ := load.Avg()
	return avg.Load1
}

func GetSystemState() (state *SystemState) {
	state = &SystemState{
		time.Now().UTC(),
		SystemMemoryConsumed(),
		SystemDiskConsumed(),
		SystemLoadAverage(),
	}
	return
}
