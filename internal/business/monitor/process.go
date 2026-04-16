package monitor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

type ProcessInfo struct {
	PID     int32
	Name    string
	CPU     float64
	Memory  float64
	Status  string
	Command string
}

type processBusiness struct{}

var ProcessBusiness = &processBusiness{}

func (b *processBusiness) ListProcesses(sortBy string, limit int) ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败: %v", err)
	}

	var processes []ProcessInfo
	for _, p := range procs {
		info, err := b.getProcessInfo(p)
		if err != nil {
			continue
		}
		processes = append(processes, info)
	}

	b.sortProcesses(processes, sortBy)

	if limit > 0 && len(processes) > limit {
		processes = processes[:limit]
	}

	return processes, nil
}

func (b *processBusiness) TopProcesses(sortBy string, limit int) ([]ProcessInfo, error) {
	if sortBy == "" {
		sortBy = "cpu"
	}
	return b.ListProcesses(sortBy, limit)
}

func (b *processBusiness) KillProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("进程不存在: %v", err)
	}

	name, _ := p.Name()

	if runtime.GOOS == "windows" {
		cmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(int(pid)))
		return cmd.Run()
	}

	cmd := exec.Command("kill", "-9", strconv.Itoa(int(pid)))
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("杀死进程失败: %v", err)
	}

	fmt.Printf("✓ 已杀死进程 [PID: %d, 名称: %s]\n", pid, name)
	return nil
}

func (b *processBusiness) FindProcess(name string) ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败: %v", err)
	}

	var processes []ProcessInfo
	nameLower := strings.ToLower(name)
	for _, p := range procs {
		pName, err := p.Name()
		if err != nil {
			continue
		}
		if strings.Contains(strings.ToLower(pName), nameLower) {
			info, err := b.getProcessInfo(p)
			if err != nil {
				continue
			}
			processes = append(processes, info)
		}
	}

	return processes, nil
}

func (b *processBusiness) getProcessInfo(p *process.Process) (ProcessInfo, error) {
	name, err := p.Name()
	if err != nil {
		return ProcessInfo{}, err
	}

	cpu, _ := p.CPUPercent()
	mem, _ := p.MemoryPercent()
	status, _ := p.Status()
	cmdline, _ := p.Cmdline()

	return ProcessInfo{
		PID:     p.Pid,
		Name:    name,
		CPU:     cpu,
		Memory:  float64(mem),
		Status:  strings.Join(status, ", "),
		Command: cmdline,
	}, nil
}

func (b *processBusiness) sortProcesses(processes []ProcessInfo, sortBy string) {
	switch sortBy {
	case "cpu":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].CPU > processes[j].CPU
		})
	case "memory":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].Memory > processes[j].Memory
		})
	case "pid":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].PID < processes[j].PID
		})
	default:
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].CPU > processes[j].CPU
		})
	}
}

func (b *processBusiness) FormatProcessTable(processes []ProcessInfo) string {
	if len(processes) == 0 {
		return "没有找到进程"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-8s %-30s %-10s %-10s %-10s\n", "PID", "名称", "CPU%", "内存%", "状态"))
	sb.WriteString(strings.Repeat("-", 80) + "\n")

	for _, p := range processes {
		name := p.Name
		if len(name) > 28 {
			name = name[:25] + "..."
		}
		sb.WriteString(fmt.Sprintf("%-8d %-30s %-10.1f %-10.1f %-10s\n",
			p.PID, name, p.CPU, p.Memory, p.Status))
	}

	return sb.String()
}

func (b *processBusiness) GetProcessCount() (int, error) {
	procs, err := process.Processes()
	if err != nil {
		return 0, err
	}
	return len(procs), nil
}

func (b *processBusiness) IsProcessRunning(pid int32) bool {
	_, err := os.FindProcess(int(pid))
	if err != nil {
		return false
	}

	p, err := process.NewProcess(pid)
	if err != nil {
		return false
	}

	running, err := p.IsRunning()
	return err == nil && running
}
