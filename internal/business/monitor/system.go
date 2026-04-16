package monitor

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

// SystemInfo 系统信息结构体
type SystemInfo struct {
	Hostname        string // 主机名
	OS              string // 操作系统
	Platform        string // 平台名称
	PlatformVersion string // 平台版本
	KernelVersion   string // 内核版本
	Architecture    string // 架构
	Uptime          uint64 // 运行时间(秒)
}

// systemBusiness 系统信息业务逻辑
type systemBusiness struct{}

// SystemBusiness 系统信息业务实例
var SystemBusiness = &systemBusiness{}

// Info 获取系统信息
func (b *systemBusiness) Info() (*SystemInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("获取系统信息失败: %v", err)
	}

	return &SystemInfo{
		Hostname:        hostInfo.Hostname,
		OS:              hostInfo.OS,
		Platform:        hostInfo.Platform,
		PlatformVersion: hostInfo.PlatformVersion,
		KernelVersion:   hostInfo.KernelVersion,
		Architecture:    hostInfo.KernelArch,
		Uptime:          hostInfo.Uptime,
	}, nil
}

// Uptime 获取系统运行时间
func (b *systemBusiness) Uptime() (uint64, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return 0, fmt.Errorf("获取运行时间失败: %v", err)
	}
	return hostInfo.Uptime, nil
}

// FormatSystemInfo 格式化系统信息为字符串
func (b *systemBusiness) FormatSystemInfo(info *SystemInfo) string {
	var sb strings.Builder

	sb.WriteString("========================================\n")
	sb.WriteString("              系统信息                  \n")
	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("主机名:       %s\n", info.Hostname))
	sb.WriteString(fmt.Sprintf("操作系统:     %s\n", info.OS))
	sb.WriteString(fmt.Sprintf("平台:         %s\n", info.Platform))
	sb.WriteString(fmt.Sprintf("平台版本:     %s\n", info.PlatformVersion))
	sb.WriteString(fmt.Sprintf("内核版本:     %s\n", info.KernelVersion))
	sb.WriteString(fmt.Sprintf("架构:         %s\n", info.Architecture))
	sb.WriteString(fmt.Sprintf("运行时间:     %s\n", b.formatUptime(info.Uptime)))
	sb.WriteString("========================================\n")

	return sb.String()
}

// FormatUptime 格式化运行时间为字符串
func (b *systemBusiness) FormatUptime(uptime uint64) string {
	return b.formatUptime(uptime)
}

// formatUptime 将秒数格式化为易读的时间字符串
func (b *systemBusiness) formatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d天", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d小时", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d分钟", minutes))
	}
	parts = append(parts, fmt.Sprintf("%d秒", secs))

	return strings.Join(parts, " ")
}

// BootTime 获取系统启动时间
func (b *systemBusiness) BootTime() (uint64, error) {
	bootTime, err := host.BootTime()
	if err != nil {
		return 0, fmt.Errorf("获取启动时间失败: %v", err)
	}
	return bootTime, nil
}

// Users 获取登录用户列表
func (b *systemBusiness) Users() ([]host.UserStat, error) {
	users, err := host.Users()
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %v", err)
	}
	return users, nil
}

// GetBootTimeFormatted 获取格式化的启动时间
func (b *systemBusiness) GetBootTimeFormatted() (string, error) {
	bootTime, err := b.BootTime()
	if err != nil {
		return "", err
	}

	t := time.Unix(int64(bootTime), 0)
	return t.Format("2006-01-02 15:04:05"), nil
}
