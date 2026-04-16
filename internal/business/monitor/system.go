package monitor

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
)

type SystemInfo struct {
	Hostname   string `json:"hostname"`
	Uptime     uint64 `json:"uptime"`
	OS         string `json:"os"`
	Platform   string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	KernelVersion string `json:"kernel_version"`
	KernelArch    string `json:"kernel_arch"`
	Procs      uint64 `json:"procs"`
}

type UptimeInfo struct {
	Uptime   uint64        `json:"uptime"`
	Duration time.Duration `json:"duration"`
	Days     int           `json:"days"`
	Hours    int           `json:"hours"`
	Minutes  int           `json:"minutes"`
	Load1    float64       `json:"load1"`
	Load5    float64       `json:"load5"`
	Load15   float64       `json:"load15"`
}

type SystemBusiness struct{}

var MonitorSystem = &SystemBusiness{}

func (b *SystemBusiness) Info(ctx context.Context) (*SystemInfo, error) {
	info, err := host.Info()
	if err != nil {
		return nil, err
	}

	return &SystemInfo{
		Hostname:        info.Hostname,
		Uptime:          info.Uptime,
		OS:              info.OS,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
		Procs:           info.Procs,
	}, nil
}

func (b *SystemBusiness) Uptime(ctx context.Context) (*UptimeInfo, error) {
	info, err := host.Info()
	if err != nil {
		return nil, err
	}

	loadInfo, _ := load.Avg()

	duration := time.Duration(info.Uptime) * time.Second
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	return &UptimeInfo{
		Uptime:   info.Uptime,
		Duration: duration,
		Days:     days,
		Hours:    hours,
		Minutes:  minutes,
		Load1:    loadInfo.Load1,
		Load5:    loadInfo.Load5,
		Load15:   loadInfo.Load15,
	}, nil
}
