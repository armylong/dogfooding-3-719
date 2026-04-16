package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"time"

	monitorBusiness "github.com/armylong/armylong-go/internal/business/monitor"
	"github.com/spf13/cobra"
)

// monitorCmd 系统监控命令结构体
type monitorCmd struct{}

// MonitorCmd 系统监控命令单例实例
var MonitorCmd = &monitorCmd{}

// MonitorHandler 系统监控命令统一入口
// 使用方式: go run main.go monitor <module> <action> [参数]
// 示例:
//
//	go run main.go monitor process list
//	go run main.go monitor process top --sort cpu --limit 10
//	go run main.go monitor disk usage
//	go run main.go monitor system info
func (m *monitorCmd) MonitorHandler(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	// 读取全局参数
	refresh, _ := cmd.Flags().GetBool("refresh")
	interval, _ := cmd.Flags().GetInt("interval")
	sortBy, _ := cmd.Flags().GetString("sort")
	limit, _ := cmd.Flags().GetInt("limit")
	perCPU, _ := cmd.Flags().GetBool("per-cpu")

	// 参数校验
	if len(args) < 1 {
		printMonitorHelp()
		return
	}

	module := args[0]
	action := ""
	if len(args) > 1 {
		action = args[1]
	}

	// 根据模块分发命令
	switch module {
	case "process":
		m.handleProcess(ctx, action, args, refresh, interval, sortBy, limit)
	case "disk":
		m.handleDisk(ctx, action, refresh, interval)
	case "memory":
		m.handleMemory(ctx, action, refresh, interval)
	case "cpu":
		m.handleCPU(ctx, action, refresh, interval, perCPU)
	case "network":
		m.handleNetwork(ctx, action, args, limit)
	case "system":
		m.handleSystem(ctx, action, refresh, interval)
	default:
		fmt.Printf("未知模块: %s\n", module)
		printMonitorHelp()
	}
}

// handleProcess 处理进程管理相关命令
func (m *monitorCmd) handleProcess(ctx context.Context, action string, args []string, refresh bool, interval int, sortBy string, limit int) {
	switch action {
	case "list":
		// 列出所有进程
		for {
			clearScreen()
			processes, err := monitorBusiness.MonitorProcess.Top(ctx, sortBy, limit)
			if err != nil {
				fmt.Printf("获取进程列表失败: %v\n", err)
				return
			}
			fmt.Println("===== 进程列表 =====")
			printProcessTable(processes)
			if !refresh {
				break
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}

	case "top":
		// 显示CPU/内存TOP进程
		if sortBy == "" {
			sortBy = "cpu"
		}
		for {
			clearScreen()
			processes, err := monitorBusiness.MonitorProcess.Top(ctx, sortBy, limit)
			if err != nil {
				fmt.Printf("获取TOP进程失败: %v\n", err)
				return
			}
			fmt.Printf("===== TOP 进程 (按%s排序, 前%d个) =====\n", sortBy, limit)
			printProcessTable(processes)
			if !refresh {
				break
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}

	case "kill":
		// 杀死指定进程
		if len(args) < 3 {
			fmt.Println("使用方式: monitor process kill <pid>")
			return
		}
		pid, err := monitorBusiness.ParsePID(args[2])
		if err != nil {
			fmt.Printf("PID格式错误: %v\n", err)
			return
		}
		err = monitorBusiness.MonitorProcess.Kill(ctx, pid)
		if err != nil {
			fmt.Printf("杀死进程失败: %v\n", err)
			return
		}
		fmt.Printf("✓ 进程 %d 已终止\n", pid)

	case "find":
		// 按名称查找进程
		if len(args) < 3 {
			fmt.Println("使用方式: monitor process find <name>")
			return
		}
		name := args[2]
		processes, err := monitorBusiness.MonitorProcess.Find(ctx, name)
		if err != nil {
			fmt.Printf("查找进程失败: %v\n", err)
			return
		}
		fmt.Printf("===== 查找进程: %s (找到 %d 个) =====\n", name, len(processes))
		printProcessTable(processes)

	default:
		fmt.Printf("未知进程命令: %s\n", action)
		fmt.Println("可用命令: list, top, kill <pid>, find <name>")
	}
}

// handleDisk 处理磁盘监控相关命令
func (m *monitorCmd) handleDisk(ctx context.Context, action string, refresh bool, interval int) {
	switch action {
	case "usage":
		// 显示磁盘使用情况
		for {
			clearScreen()
			usage, err := monitorBusiness.MonitorDisk.Usage(ctx, "/")
			if err != nil {
				fmt.Printf("获取磁盘使用情况失败: %v\n", err)
				return
			}
			fmt.Println("===== 磁盘使用情况 =====")
			fmt.Printf("路径: %s\n", usage.Path)
			fmt.Printf("文件系统: %s\n", usage.Fstype)
			fmt.Printf("总容量: %s\n", formatBytes(usage.Total))
			fmt.Printf("已使用: %s (%.2f%%)\n", formatBytes(usage.Used), usage.UsedPercent)
			fmt.Printf("可用: %s\n", formatBytes(usage.Free))
			fmt.Print("使用率: [")
			used := int(usage.UsedPercent / 5)
			for i := 0; i < 20; i++ {
				if i < used {
					fmt.Print("█")
				} else {
					fmt.Print("░")
				}
			}
			fmt.Printf("] %.1f%%\n", usage.UsedPercent)
			if !refresh {
				break
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}

	case "list":
		// 列出所有磁盘分区
		partitions, err := monitorBusiness.MonitorDisk.List(ctx)
		if err != nil {
			fmt.Printf("获取磁盘分区失败: %v\n", err)
			return
		}
		fmt.Println("===== 磁盘分区列表 =====")
		fmt.Printf("%-30s %-30s %-15s %s\n", "设备", "挂载点", "文件系统", "选项")
		fmt.Println("--------------------------------------------------------------------------------")
		for _, p := range partitions {
			opts := ""
			if len(p.Opts) > 0 {
				opts = p.Opts[0]
				for i := 1; i < len(p.Opts); i++ {
					opts += "," + p.Opts[i]
				}
			}
			fmt.Printf("%-30s %-30s %-15s %s\n", p.Device, p.Mountpoint, p.Fstype, opts)
		}

	default:
		fmt.Printf("未知磁盘命令: %s\n", action)
		fmt.Println("可用命令: usage, list")
	}
}

// handleMemory 处理内存监控相关命令
func (m *monitorCmd) handleMemory(ctx context.Context, action string, refresh bool, interval int) {
	switch action {
	case "usage":
		// 显示内存使用情况
		for {
			clearScreen()
			mem, err := monitorBusiness.MonitorMemory.Usage(ctx)
			if err != nil {
				fmt.Printf("获取内存使用情况失败: %v\n", err)
				return
			}
			fmt.Println("===== 内存使用情况 =====")
			fmt.Printf("总内存: %s\n", formatBytes(mem.Total))
			fmt.Printf("已使用: %s (%.2f%%)\n", formatBytes(mem.Used), mem.UsedPercent)
			fmt.Printf("可用: %s\n", formatBytes(mem.Available))
			fmt.Printf("空闲: %s\n", formatBytes(mem.Free))
			fmt.Printf("活跃: %s\n", formatBytes(mem.Active))
			fmt.Printf("缓存: %s\n", formatBytes(mem.Cached))
			fmt.Print("使用率: [")
			used := int(mem.UsedPercent / 5)
			for i := 0; i < 20; i++ {
				if i < used {
					fmt.Print("█")
				} else {
					fmt.Print("░")
				}
			}
			fmt.Printf("] %.1f%%\n", mem.UsedPercent)
			if !refresh {
				break
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}

	default:
		fmt.Printf("未知内存命令: %s\n", action)
		fmt.Println("可用命令: usage")
	}
}

// handleCPU 处理CPU监控相关命令
func (m *monitorCmd) handleCPU(ctx context.Context, action string, refresh bool, interval int, perCPU bool) {
	switch action {
	case "usage":
		// 显示CPU使用率
		for {
			clearScreen()
			usage, err := monitorBusiness.MonitorCPU.Usage(ctx, perCPU)
			if err != nil {
				fmt.Printf("获取CPU使用率失败: %v\n", err)
				return
			}
			fmt.Println("===== CPU使用率 =====")
			fmt.Print("平均使用率: [")
			used := int(usage.Average / 5)
			for i := 0; i < 20; i++ {
				if i < used {
					fmt.Print("█")
				} else {
					fmt.Print("░")
				}
			}
			fmt.Printf("] %.1f%%\n", usage.Average)
			if perCPU {
				for i, p := range usage.Percent {
					fmt.Printf("CPU%d: [", i)
					cpuUsed := int(p / 5)
					for j := 0; j < 20; j++ {
						if j < cpuUsed {
							fmt.Print("█")
						} else {
							fmt.Print("░")
						}
					}
					fmt.Printf("] %.1f%%\n", p)
				}
			}
			if !refresh {
				break
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}

	case "info":
		// 显示CPU信息
		infos, err := monitorBusiness.MonitorCPU.Info(ctx)
		if err != nil {
			fmt.Printf("获取CPU信息失败: %v\n", err)
			return
		}
		fmt.Println("===== CPU信息 =====")
		for _, info := range infos {
			fmt.Printf("处理器 CPU %d:\n", info.CPU)
			fmt.Printf("  型号: %s\n", info.ModelName)
			fmt.Printf("  厂商: %s\n", info.VendorID)
			fmt.Printf("  主频: %.2f MHz\n", info.Mhz)
			fmt.Printf("  核心数: %d\n", info.Cores)
			fmt.Printf("  缓存: %d KB\n", info.CacheSize)
			fmt.Println()
		}

	default:
		fmt.Printf("未知CPU命令: %s\n", action)
		fmt.Println("可用命令: usage, info")
	}
}

// handleNetwork 处理网络监控相关命令
func (m *monitorCmd) handleNetwork(ctx context.Context, action string, args []string, limit int) {
	switch action {
	case "connections":
		// 显示网络连接
		conns, err := monitorBusiness.MonitorNetwork.Connections(ctx)
		if err != nil {
			fmt.Printf("获取网络连接失败: %v\n", err)
			return
		}
		fmt.Println("===== 网络连接 =====")
		fmt.Printf("%-6s %-22s %-22s %-15s %-6s %s\n", "协议", "本地地址", "远程地址", "状态", "PID", "进程名")
		fmt.Println("----------------------------------------------------------------------------------------")
		count := 0
		for _, c := range conns {
			if limit > 0 && count >= limit {
				break
			}
			if c.Status == "" || c.Status == "NONE" {
				continue
			}
			proto := "TCP"
			if c.Type == 2 {
				proto = "UDP"
			}
			name := ""
			if c.PID > 0 {
				p, _ := monitorBusiness.ProcessNewProcess(c.PID)
				if p != nil {
					name, _ = p.Name()
				}
			}
			fmt.Printf("%-6s %-22s %-22s %-15s %-6d %s\n", proto, c.Laddr, c.Raddr, c.Status, c.PID, name)
			count++
		}

	case "ports":
		// 显示端口占用
		ports, err := monitorBusiness.MonitorNetwork.Ports(ctx)
		if err != nil {
			fmt.Printf("获取端口占用失败: %v\n", err)
			return
		}
		sort.Slice(ports, func(i, j int) bool {
			return ports[i].Port < ports[j].Port
		})
		fmt.Println("===== 端口占用 =====")
		fmt.Printf("%-10s %-10s %-15s %-8s %s\n", "端口", "协议", "状态", "PID", "进程名")
		fmt.Println("----------------------------------------------------------------")
		for _, p := range ports {
			fmt.Printf("%-10d %-10s %-15s %-8d %s\n", p.Port, p.Protocol, p.Status, p.PID, p.Name)
		}

	case "kill-port":
		// 杀死占用指定端口的进程
		if len(args) < 3 {
			fmt.Println("使用方式: monitor network kill-port <port>")
			return
		}
		port, err := monitorBusiness.ParsePort(args[2])
		if err != nil {
			fmt.Printf("端口格式错误: %v\n", err)
			return
		}
		err = monitorBusiness.MonitorNetwork.KillPort(ctx, port)
		if err != nil {
			fmt.Printf("杀死端口进程失败: %v\n", err)
			return
		}
		fmt.Printf("✓ 端口 %d 相关进程已终止\n", port)

	default:
		fmt.Printf("未知网络命令: %s\n", action)
		fmt.Println("可用命令: connections, ports, kill-port <port>")
	}
}

// handleSystem 处理系统信息相关命令
func (m *monitorCmd) handleSystem(ctx context.Context, action string, refresh bool, interval int) {
	switch action {
	case "info":
		// 显示系统信息
		info, err := monitorBusiness.MonitorSystem.Info(ctx)
		if err != nil {
			fmt.Printf("获取系统信息失败: %v\n", err)
			return
		}
		fmt.Println("===== 系统信息 =====")
		fmt.Printf("主机名: %s\n", info.Hostname)
		fmt.Printf("操作系统: %s\n", info.OS)
		fmt.Printf("平台: %s %s\n", info.Platform, info.PlatformVersion)
		fmt.Printf("内核版本: %s\n", info.KernelVersion)
		fmt.Printf("架构: %s\n", info.KernelArch)
		fmt.Printf("进程数: %d\n", info.Procs)

	case "uptime":
		// 显示系统运行时间
		for {
			clearScreen()
			uptime, err := monitorBusiness.MonitorSystem.Uptime(ctx)
			if err != nil {
				fmt.Printf("获取系统运行时间失败: %v\n", err)
				return
			}
			fmt.Println("===== 系统运行时间 =====")
			fmt.Printf("运行时间: %d天 %d小时 %d分钟\n", uptime.Days, uptime.Hours, uptime.Minutes)
			fmt.Printf("总秒数: %d\n", uptime.Uptime)
			fmt.Printf("负载平均: 1分钟=%.2f, 5分钟=%.2f, 15分钟=%.2f\n", uptime.Load1, uptime.Load5, uptime.Load15)
			if !refresh {
				break
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}

	default:
		fmt.Printf("未知系统命令: %s\n", action)
		fmt.Println("可用命令: info, uptime")
	}
}

// printMonitorHelp 打印监控命令帮助信息
func printMonitorHelp() {
	fmt.Println("系统监控命令使用方式:")
	fmt.Println()
	fmt.Println("  进程管理:")
	fmt.Println("    monitor process list                    列出所有进程")
	fmt.Println("    monitor process top                     显示TOP进程")
	fmt.Println("    monitor process kill <pid>              杀死指定进程")
	fmt.Println("    monitor process find <name>             按名称查找进程")
	fmt.Println()
	fmt.Println("  磁盘监控:")
	fmt.Println("    monitor disk usage                      显示磁盘使用情况")
	fmt.Println("    monitor disk list                       列出所有磁盘分区")
	fmt.Println()
	fmt.Println("  内存监控:")
	fmt.Println("    monitor memory usage                    显示内存使用情况")
	fmt.Println()
	fmt.Println("  CPU监控:")
	fmt.Println("    monitor cpu usage                       显示CPU使用率")
	fmt.Println("    monitor cpu info                        显示CPU信息")
	fmt.Println()
	fmt.Println("  网络监控:")
	fmt.Println("    monitor network connections             显示网络连接")
	fmt.Println("    monitor network ports                   显示端口占用")
	fmt.Println("    monitor network kill-port <port>        杀死占用指定端口的进程")
	fmt.Println()
	fmt.Println("  系统信息:")
	fmt.Println("    monitor system info                     显示系统信息")
	fmt.Println("    monitor system uptime                   显示系统运行时间")
	fmt.Println()
	fmt.Println("  全局参数:")
	fmt.Println("    --refresh                               实时刷新显示")
	fmt.Println("    --interval <秒>                         刷新间隔，默认2秒")
	fmt.Println("    --sort <cpu/memory/pid>                 排序方式")
	fmt.Println("    --limit <数量>                          显示数量限制，默认10条")
	fmt.Println("    --per-cpu                               显示每个CPU核心使用率")
}

// printProcessTable 格式化打印进程表格
func printProcessTable(processes []monitorBusiness.ProcessInfo) {
	fmt.Printf("%-8s %-25s %-10s %-10s %-15s %-15s %s\n", "PID", "进程名", "CPU%", "内存%", "内存(RSS)", "状态", "用户")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, p := range processes {
		fmt.Printf("%-8d %-25s %-10.1f %-10.1f %-15s %-15s %s\n",
			p.PID,
			truncateString(p.Name, 25),
			p.CPUPercent,
			p.MemPercent,
			formatBytes(p.MemRSS),
			p.Status,
			p.Username)
	}
}

// formatBytes 格式化字节大小为可读字符串
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// truncateString 截断字符串到指定长度
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// clearScreen 清空终端屏幕，跨平台支持
func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}
