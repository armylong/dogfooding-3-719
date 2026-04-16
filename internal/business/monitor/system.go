package monitor

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
)

type SystemInfo struct {
	Hostname        string
	OS              string
	Platform        string
	PlatformVersion string
	KernelVersion   string
	Architecture    string
	Uptime          uint64
	BootTime        uint64
	Procs           uint64
}

type UptimeInfo struct {
	Days    int
	Hours   int
	Minutes int
	Seconds int
	Total   uint64
}

type LoadInfo struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

type systemBusiness struct{}

var SystemBusiness = &systemBusiness{}

func (b *systemBusiness) Info() (*SystemInfo, error) {
	info, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("获取系统信息失败: %v", err)
	}

	return &SystemInfo{
		Hostname:        info.Hostname,
		OS:              info.OS,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		Architecture:    runtime.GOARCH,
		Uptime:          info.Uptime,
		BootTime:        info.BootTime,
		Procs:           info.Procs,
	}, nil
}

func (b *systemBusiness) Uptime() (*UptimeInfo, error) {
	info, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("获取运行时间失败: %v", err)
	}

	uptime := info.Uptime
	return &UptimeInfo{
		Days:    int(uptime / 86400),
		Hours:   int((uptime % 86400) / 3600),
		Minutes: int((uptime % 3600) / 60),
		Seconds: int(uptime % 60),
		Total:   uptime,
	}, nil
}

func (b *systemBusiness) Load() (*LoadInfo, error) {
	loadStat, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("获取系统负载失败: %v", err)
	}

	return &LoadInfo{
		Load1:  loadStat.Load1,
		Load5:  loadStat.Load5,
		Load15: loadStat.Load15,
	}, nil
}

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
	sb.WriteString(fmt.Sprintf("进程数:       %d\n", info.Procs))
	sb.WriteString("========================================\n")

	return sb.String()
}

func (b *systemBusiness) FormatUptime(uptime *UptimeInfo) string {
	var sb strings.Builder

	sb.WriteString("========================================\n")
	sb.WriteString("            系统运行时间                \n")
	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("运行天数:     %d 天\n", uptime.Days))
	sb.WriteString(fmt.Sprintf("运行时间:     %d 小时 %d 分钟 %d 秒\n", uptime.Hours, uptime.Minutes, uptime.Seconds))
	sb.WriteString(fmt.Sprintf("总秒数:       %d 秒\n", uptime.Total))
	sb.WriteString("========================================\n")

	return sb.String()
}

func (b *systemBusiness) FormatLoad(load *LoadInfo) string {
	var sb strings.Builder

	sb.WriteString("========================================\n")
	sb.WriteString("            系统负载                    \n")
	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("1分钟负载:    %.2f\n", load.Load1))
	sb.WriteString(fmt.Sprintf("5分钟负载:    %.2f\n", load.Load5))
	sb.WriteString(fmt.Sprintf("15分钟负载:   %.2f\n", load.Load15))
	sb.WriteString("========================================\n")

	return sb.String()
}

func (b *systemBusiness) GetBootTime() (time.Time, error) {
	bootTime, err := host.BootTime()
	if err != nil {
		return time.Time{}, fmt.Errorf("获取启动时间失败: %v", err)
	}
	return time.Unix(int64(bootTime), 0), nil
}

func (b *systemBusiness) GetUsers() ([]host.UserStat, error) {
	users, err := host.Users()
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}
	return users, nil
}
