// Package monitor 系统监控业务层
// 提供进程、磁盘、内存、CPU、网络、系统信息等监控功能
package monitor

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo 进程信息结构体
type ProcessInfo struct {
	PID        int32   `json:"pid"`        // 进程ID
	Name       string  `json:"name"`       // 进程名称
	CPUPercent float64 `json:"cpu_percent"` // CPU使用率
	MemPercent float32 `json:"mem_percent"` // 内存使用率
	MemRSS     uint64  `json:"mem_rss"`    // 内存占用RSS
	Status     string  `json:"status"`     // 进程状态
	Username   string  `json:"username"`   // 运行用户
}

// ProcessBusiness 进程管理业务结构体
type ProcessBusiness struct{}

// MonitorProcess 进程管理业务单例
var MonitorProcess = &ProcessBusiness{}

// List 获取所有进程列表
func (b *ProcessBusiness) List(ctx context.Context) ([]ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var result []ProcessInfo
	for _, p := range processes {
		info := b.getProcessInfo(p)
		result = append(result, info)
	}

	return result, nil
}

// Top 获取排序后的进程列表，支持按cpu/memory/pid排序
func (b *ProcessBusiness) Top(ctx context.Context, sortBy string, limit int) ([]ProcessInfo, error) {
	processes, err := b.List(ctx)
	if err != nil {
		return nil, err
	}

	// 根据排序方式排序
	switch sortBy {
	case "cpu":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].CPUPercent > processes[j].CPUPercent
		})
	case "memory":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].MemPercent > processes[j].MemPercent
		})
	case "pid":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].PID > processes[j].PID
		})
	default:
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].CPUPercent > processes[j].CPUPercent
		})
	}

	// 限制返回数量
	if limit > 0 && limit < len(processes) {
		return processes[:limit], nil
	}
	return processes, nil
}

// Kill 杀死指定PID的进程
func (b *ProcessBusiness) Kill(ctx context.Context, pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}

// Find 按名称模糊查找进程
func (b *ProcessBusiness) Find(ctx context.Context, name string) ([]ProcessInfo, error) {
	processes, err := b.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []ProcessInfo
	for _, p := range processes {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(name)) {
			result = append(result, p)
		}
	}

	return result, nil
}

// getProcessInfo 获取单个进程的详细信息
func (b *ProcessBusiness) getProcessInfo(p *process.Process) ProcessInfo {
	info := ProcessInfo{PID: p.Pid}

	if name, err := p.Name(); err == nil {
		info.Name = name
	}

	if cpu, err := p.CPUPercent(); err == nil {
		info.CPUPercent = cpu
	}

	if mem, err := p.MemoryPercent(); err == nil {
		info.MemPercent = mem
	}

	if memInfo, err := p.MemoryInfo(); err == nil && memInfo != nil {
		info.MemRSS = memInfo.RSS
	}

	if status, err := p.Status(); err == nil && len(status) > 0 {
		info.Status = status[0]
	}

	if username, err := p.Username(); err == nil {
		info.Username = username
	}

	return info
}

// ParsePID 解析字符串PID为int32
func ParsePID(pidStr string) (int32, error) {
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, err
	}
	return int32(pid), nil
}

// ProcessNewProcess 创建进程实例
func ProcessNewProcess(pid int32) (*process.Process, error) {
	return process.NewProcess(pid)
}
