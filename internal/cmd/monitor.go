package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	monitorBusiness "github.com/armylong/armylong-go/internal/business/monitor"
	"github.com/spf13/cobra"
)

type monitorCmd struct{}

var MonitorCmd = &monitorCmd{}

var (
	RefreshFlag  bool
	IntervalFlag int
	SortFlag     string
	LimitFlag    int
)

func (m *monitorCmd) MonitorHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("错误: 请指定子命令")
		fmt.Println("可用命令: process, disk, memory, cpu, network, system")
		return
	}

	switch args[0] {
	case "process":
		m.handleProcess(cmd, args[1:])
	case "disk":
		m.handleDisk(cmd, args[1:])
	case "memory":
		m.handleMemory(cmd, args[1:])
	case "cpu":
		m.handleCPU(cmd, args[1:])
	case "network":
		m.handleNetwork(cmd, args[1:])
	case "system":
		m.handleSystem(cmd, args[1:])
	default:
		fmt.Printf("未知命令: %s\n", args[0])
		fmt.Println("可用命令: process, disk, memory, cpu, network, system")
	}
}

func (m *monitorCmd) handleProcess(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("错误: 请指定 process 子命令")
		fmt.Println("可用命令: list, top, kill, find")
		return
	}

	switch args[0] {
	case "list":
		m.processList()
	case "top":
		m.processTop()
	case "kill":
		if len(args) < 2 {
			fmt.Println("错误: 请指定进程PID")
			return
		}
		m.processKill(args[1])
	case "find":
		if len(args) < 2 {
			fmt.Println("错误: 请指定进程名称")
			return
		}
		m.processFind(args[1])
	default:
		fmt.Printf("未知命令: %s\n", args[0])
		fmt.Println("可用命令: list, top, kill, find")
	}
}

func (m *monitorCmd) processList() {
	if RefreshFlag {
		m.runWithRefresh(func() {
			processes, err := monitorBusiness.ProcessBusiness.ListProcesses(SortFlag, LimitFlag)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
				return
			}
			fmt.Println(monitorBusiness.ProcessBusiness.FormatProcessTable(processes))
		})
	} else {
		processes, err := monitorBusiness.ProcessBusiness.ListProcesses(SortFlag, LimitFlag)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		fmt.Println(monitorBusiness.ProcessBusiness.FormatProcessTable(processes))
	}
}

func (m *monitorCmd) processTop() {
	if RefreshFlag {
		m.runWithRefresh(func() {
			processes, err := monitorBusiness.ProcessBusiness.TopProcesses(SortFlag, LimitFlag)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
				return
			}
			fmt.Println(monitorBusiness.ProcessBusiness.FormatProcessTable(processes))
		})
	} else {
		processes, err := monitorBusiness.ProcessBusiness.TopProcesses(SortFlag, LimitFlag)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		fmt.Println(monitorBusiness.ProcessBusiness.FormatProcessTable(processes))
	}
}

func (m *monitorCmd) processKill(pidStr string) {
	pid, err := strconv.ParseInt(pidStr, 10, 32)
	if err != nil {
		fmt.Printf("错误: 无效的PID: %v\n", err)
		return
	}

	err = monitorBusiness.ProcessBusiness.KillProcess(int32(pid))
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	}
}

func (m *monitorCmd) processFind(name string) {
	processes, err := monitorBusiness.ProcessBusiness.FindProcess(name)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	if len(processes) == 0 {
		fmt.Printf("未找到匹配 '%s' 的进程\n", name)
		return
	}

	fmt.Printf("找到 %d 个匹配 '%s' 的进程:\n", len(processes), name)
	fmt.Println(monitorBusiness.ProcessBusiness.FormatProcessTable(processes))
}

func (m *monitorCmd) handleDisk(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("错误: 请指定 disk 子命令")
		fmt.Println("可用命令: usage, list")
		return
	}

	switch args[0] {
	case "usage":
		m.diskUsage()
	case "list":
		m.diskList()
	default:
		fmt.Printf("未知命令: %s\n", args[0])
		fmt.Println("可用命令: usage, list")
	}
}

func (m *monitorCmd) diskUsage() {
	if RefreshFlag {
		m.runWithRefresh(func() {
			disks, err := monitorBusiness.DiskBusiness.Usage()
			if err != nil {
				fmt.Printf("错误: %v\n", err)
				return
			}
			fmt.Println(monitorBusiness.DiskBusiness.FormatDiskTable(disks))
		})
	} else {
		disks, err := monitorBusiness.DiskBusiness.Usage()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		fmt.Println(monitorBusiness.DiskBusiness.FormatDiskTable(disks))
	}
}

func (m *monitorCmd) diskList() {
	disks, err := monitorBusiness.DiskBusiness.List()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Println(monitorBusiness.DiskBusiness.FormatDiskTable(disks))
}

func (m *monitorCmd) handleMemory(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("错误: 请指定 memory 子命令")
		fmt.Println("可用命令: usage")
		return
	}

	switch args[0] {
	case "usage":
		m.memoryUsage()
	default:
		fmt.Printf("未知命令: %s\n", args[0])
		fmt.Println("可用命令: usage")
	}
}

func (m *monitorCmd) memoryUsage() {
	if RefreshFlag {
		m.runWithRefresh(func() {
			info, err := monitorBusiness.MemoryBusiness.Usage()
			if err != nil {
				fmt.Printf("错误: %v\n", err)
				return
			}
			fmt.Println(monitorBusiness.MemoryBusiness.FormatMemoryTable(info))
		})
	} else {
		info, err := monitorBusiness.MemoryBusiness.Usage()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		fmt.Println(monitorBusiness.MemoryBusiness.FormatMemoryTable(info))
	}
}

func (m *monitorCmd) handleCPU(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("错误: 请指定 cpu 子命令")
		fmt.Println("可用命令: usage, info")
		return
	}

	switch args[0] {
	case "usage":
		m.cpuUsage()
	case "info":
		m.cpuInfo()
	default:
		fmt.Printf("未知命令: %s\n", args[0])
		fmt.Println("可用命令: usage, info")
	}
}

func (m *monitorCmd) cpuUsage() {
	if RefreshFlag {
		m.runWithRefresh(func() {
			usage, err := monitorBusiness.CPUBusiness.Usage(time.Second)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
				return
			}
			fmt.Println(monitorBusiness.CPUBusiness.FormatUsageOnly(usage))
		})
	} else {
		usage, err := monitorBusiness.CPUBusiness.Usage(time.Second)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		fmt.Println(monitorBusiness.CPUBusiness.FormatUsageOnly(usage))
	}
}

func (m *monitorCmd) cpuInfo() {
	usage, err := monitorBusiness.CPUBusiness.Usage(time.Second)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	info, err := monitorBusiness.CPUBusiness.Info()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	fmt.Println(monitorBusiness.CPUBusiness.FormatCPUTable(usage, info))
}

func (m *monitorCmd) handleNetwork(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("错误: 请指定 network 子命令")
		fmt.Println("可用命令: connections, ports, kill-port")
		return
	}

	switch args[0] {
	case "connections":
		m.networkConnections()
	case "ports":
		m.networkPorts()
	case "kill-port":
		if len(args) < 2 {
			fmt.Println("错误: 请指定端口号")
			return
		}
		m.networkKillPort(args[1])
	default:
		fmt.Printf("未知命令: %s\n", args[0])
		fmt.Println("可用命令: connections, ports, kill-port")
	}
}

func (m *monitorCmd) networkConnections() {
	if RefreshFlag {
		m.runWithRefresh(func() {
			connections, err := monitorBusiness.NetworkBusiness.Connections()
			if err != nil {
				fmt.Printf("错误: %v\n", err)
				return
			}
			fmt.Println(monitorBusiness.NetworkBusiness.FormatConnectionTable(connections))
		})
	} else {
		connections, err := monitorBusiness.NetworkBusiness.Connections()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		fmt.Println(monitorBusiness.NetworkBusiness.FormatConnectionTable(connections))
	}
}

func (m *monitorCmd) networkPorts() {
	if RefreshFlag {
		m.runWithRefresh(func() {
			ports, err := monitorBusiness.NetworkBusiness.Ports()
			if err != nil {
				fmt.Printf("错误: %v\n", err)
				return
			}
			fmt.Println(monitorBusiness.NetworkBusiness.FormatPortTable(ports))
		})
	} else {
		ports, err := monitorBusiness.NetworkBusiness.Ports()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		fmt.Println(monitorBusiness.NetworkBusiness.FormatPortTable(ports))
	}
}

func (m *monitorCmd) networkKillPort(portStr string) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Printf("错误: 无效的端口号: %v\n", err)
		return
	}

	err = monitorBusiness.NetworkBusiness.KillPort(port)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	}
}

func (m *monitorCmd) handleSystem(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("错误: 请指定 system 子命令")
		fmt.Println("可用命令: info, uptime")
		return
	}

	switch args[0] {
	case "info":
		m.systemInfo()
	case "uptime":
		m.systemUptime()
	default:
		fmt.Printf("未知命令: %s\n", args[0])
		fmt.Println("可用命令: info, uptime")
	}
}

func (m *monitorCmd) systemInfo() {
	info, err := monitorBusiness.SystemBusiness.Info()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Println(monitorBusiness.SystemBusiness.FormatSystemInfo(info))
}

func (m *monitorCmd) systemUptime() {
	uptime, err := monitorBusiness.SystemBusiness.Uptime()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Println(monitorBusiness.SystemBusiness.FormatUptime(uptime))
}

func (m *monitorCmd) runWithRefresh(displayFunc func()) {
	interval := time.Duration(IntervalFlag) * time.Second
	if interval < time.Second {
		interval = 2 * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	displayFunc()

	for {
		select {
		case <-sigChan:
			fmt.Println("\n停止刷新")
			return
		case <-ticker.C:
			m.clearScreen()
			displayFunc()
		}
	}
}

func (m *monitorCmd) clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func init() {
	SortFlag = "cpu"
	LimitFlag = 10
	IntervalFlag = 2
	RefreshFlag = false
}

func parseSortFlag(sortBy string) string {
	switch strings.ToLower(sortBy) {
	case "cpu", "memory", "pid":
		return strings.ToLower(sortBy)
	default:
		return "cpu"
	}
}
