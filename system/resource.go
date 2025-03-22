package system

import (
	"fmt"

	"github.com/farseer-go/fs/net"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type Resource struct {
	OS                 string  // 操作系统名称
	HostName           string  // 主机名
	Architecture       string  // 系统架构
	Processes          uint64  // 进程数
	IP                 string  // IP
	CpuName            string  // CPU名称
	CpuMhz             float64 // CPU总赫兹
	CpuCores           int     // CPU核心数
	CpuUsagePercent    float64 // CPU使用百分比
	MemoryTotal        uint64  // 总内存
	MemoryAvailable    uint64  // 内存可用量（B）
	MemoryUsage        uint64  // 内存已使用（B）
	MemoryUsagePercent float64 // 内存使用百分比
	DiskTotal          uint64  // 硬盘总容量
	DiskAvailable      uint64  // 硬盘可用空间
	DiskUsage          uint64  // 硬盘已用空间
	DiskUsagePercent   float64 // 硬盘使用百分比
}

func (receiver *Resource) ToString() string {
	return fmt.Sprintf("%+v", receiver)
}

// GetResource 获取当前环境信息
func GetResource() Resource {
	cpuPercents, _ := cpu.Percent(0, false)
	infoStats, _ := cpu.Info()
	memory, _ := mem.VirtualMemory()
	hostInfo, _ := host.Info()
	diskUsage, _ := disk.Usage("/")
	if len(infoStats) == 0 {
		infoStats = []cpu.InfoStat{{}}
	}
	if len(cpuPercents) == 0 {
		cpuPercents = []float64{float64(0)}
	}

	// 取所有核心数
	var cores int
	var cpuMhz float64
	var cpuModelName string
	for _, cpuStat := range infoStats {
		cores += int(cpuStat.Cores)
		if cpuStat.Mhz > cpuMhz {
			cpuMhz = cpuStat.Mhz
			cpuModelName = cpuStat.ModelName
		}
	}

	// cpu使用率
	var cpuUsagePercent float64
	for _, singleInfo := range cpuPercents {
		cpuUsagePercent += singleInfo
	}

	cpuUsagePercent = cpuUsagePercent / float64(len(cpuPercents))
	return Resource{
		Architecture:       hostInfo.KernelArch,
		HostName:           hostInfo.Hostname,
		OS:                 hostInfo.OS,
		Processes:          hostInfo.Procs,
		IP:                 net.GetIp(),
		CpuName:            cpuModelName,
		CpuMhz:             cpuMhz,
		CpuUsagePercent:    cpuUsagePercent,
		CpuCores:           cores,
		MemoryUsage:        memory.Used,
		MemoryUsagePercent: memory.UsedPercent,
		MemoryTotal:        memory.Total,
		MemoryAvailable:    memory.Available,
		DiskTotal:          diskUsage.Total,
		DiskAvailable:      diskUsage.Free,
		DiskUsage:          diskUsage.Used,
		DiskUsagePercent:   diskUsage.UsedPercent,
	}
}
