package monitor

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"
)

// DiskInfo 磁盘信息结构体
type DiskInfo struct {
	Device      string  // 设备名称
	MountPoint  string  // 挂载点
	Total       uint64  // 总容量(字节)
	Used        uint64  // 已使用(字节)
	Free        uint64  // 剩余空间(字节)
	UsedPercent float64 // 使用率(%)
	Fstype      string  // 文件系统类型
}

// diskBusiness 磁盘管理业务逻辑
type diskBusiness struct{}

// DiskBusiness 磁盘管理业务实例
var DiskBusiness = &diskBusiness{}

// Usage 获取磁盘使用情况
func (b *diskBusiness) Usage() ([]DiskInfo, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, fmt.Errorf("获取磁盘分区失败: %v", err)
	}

	var diskInfos []DiskInfo
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		diskInfos = append(diskInfos, DiskInfo{
			Device:      partition.Device,
			MountPoint:  partition.Mountpoint,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
			Fstype:      partition.Fstype,
		})
	}

	return diskInfos, nil
}

// List 获取磁盘分区列表
func (b *diskBusiness) List() ([]DiskInfo, error) {
	return b.Usage()
}

// FormatDiskTable 格式化磁盘信息为表格字符串
func (b *diskBusiness) FormatDiskTable(disks []DiskInfo) string {
	if len(disks) == 0 {
		return "没有找到磁盘分区"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-15s %-20s %-15s %-15s %-10s\n", "设备", "挂载点", "总容量", "已使用", "使用率"))
	sb.WriteString(strings.Repeat("-", 85) + "\n")

	for _, d := range disks {
		total := b.formatBytes(d.Total)
		used := b.formatBytes(d.Used)
		sb.WriteString(fmt.Sprintf("%-15s %-20s %-15s %-15s %-10.1f%%\n",
			d.Device, d.MountPoint, total, used, d.UsedPercent))
	}

	return sb.String()
}

// formatBytes 将字节数格式化为易读的字符串
func (b *diskBusiness) formatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2fTB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2fGB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2fMB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2fKB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

// GetRootDiskUsage 获取根分区磁盘使用情况
func (b *diskBusiness) GetRootDiskUsage() (*DiskInfo, error) {
	disks, err := b.Usage()
	if err != nil {
		return nil, err
	}

	for _, d := range disks {
		if d.MountPoint == "/" {
			return &d, nil
		}
	}

	if len(disks) > 0 {
		return &disks[0], nil
	}

	return nil, fmt.Errorf("未找到根分区")
}

// GetIOCounters 获取磁盘IO统计信息
func (b *diskBusiness) GetIOCounters() (map[string]disk.IOCountersStat, error) {
	ioCounters, err := disk.IOCounters()
	if err != nil {
		return nil, fmt.Errorf("获取磁盘IO统计失败: %v", err)
	}
	return ioCounters, nil
}
