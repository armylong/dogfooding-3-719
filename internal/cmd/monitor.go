// Package cmd 提供命令行命令实现
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

// monitorCmd 监控命令结构体
type monitorCmd struct{}

// MonitorCmd 监控命令实例
var MonitorCmd = &monitorCmd{}

// 命令行标志变量
var (
	RefreshFlag  bool   // 是否实时刷新
	IntervalFlag int    // 刷新间隔(秒)
	SortFlag     string // 排序方式(cpu/memory/pid)
	LimitFlag    int    // 显示数量限制
)

// MonitorHandler 监控命令主处理函数
// 根据子命令类型分发到对应的处理函数
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

// handleProcess 处理进程相关子命令
// 支持的子命令: list, top, kill, find
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

// processList 显示进程列表
// 支持实时刷新和排序
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

// processTop 显示资源占用最高的进程
// 支持实时刷新和排序
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

// processKill 杀死指定进程
// pidStr: 进程ID字符串
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

// processFind 按名称查找进程
// name: 进程名称(支持模糊匹配)
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

// handleDisk 处理磁盘相关子命令
// 支持的子命令: usage, list
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

// diskUsage 显示磁盘使用情况
// 支持实时刷新
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

// diskList 显示磁盘分区列表
func (m *monitorCmd) diskList() {
	disks, err := monitorBusiness.DiskBusiness.List()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Println(monitorBusiness.DiskBusiness.FormatDiskTable(disks))
}

// handleMemory 处理内存相关子命令
// 支持的子命令: usage
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

// memoryUsage 显示内存使用情况
// 支持实时刷新
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

// handleCPU 处理CPU相关子命令
// 支持的子命令: usage, info
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

// cpuUsage 显示CPU使用率
// 支持实时刷新
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

// cpuInfo 显示CPU详细信息
// 包括型号、核心数和使用率
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

// handleNetwork 处理网络相关子命令
// 支持的子命令: connections, ports, kill-port
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

// networkConnections 显示网络连接列表
// 支持实时刷新
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

// networkPorts 显示监听端口列表
// 支持实时刷新
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

// networkKillPort 杀死占用指定端口的进程
// portStr: 端口号字符串
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

// handleSystem 处理系统信息相关子命令
// 支持的子命令: info, uptime
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

// systemInfo 显示系统信息
// 包括主机名、操作系统、平台等
func (m *monitorCmd) systemInfo() {
	info, err := monitorBusiness.SystemBusiness.Info()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Println(monitorBusiness.SystemBusiness.FormatSystemInfo(info))
}

// systemUptime 显示系统运行时间
func (m *monitorCmd) systemUptime() {
	uptime, err := monitorBusiness.SystemBusiness.Uptime()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Println(monitorBusiness.SystemBusiness.FormatUptime(uptime))
}

// runWithRefresh 以实时刷新模式运行显示函数
// displayFunc: 显示内容的函数
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

// clearScreen 清屏
func (m *monitorCmd) clearScreen() {
	fmt.Print("\033[2J\033[H")
}

// init 初始化默认值
func init() {
	SortFlag = "cpu"
	LimitFlag = 10
	IntervalFlag = 2
	RefreshFlag = false
}

// parseSortFlag 解析排序参数
// 返回有效的排序方式
func parseSortFlag(sortBy string) string {
	switch strings.ToLower(sortBy) {
	case "cpu", "memory", "pid":
		return strings.ToLower(sortBy)
	default:
		return "cpu"
	}
}
