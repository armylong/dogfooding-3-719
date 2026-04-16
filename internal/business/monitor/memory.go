package monitor

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryInfo struct {
	Total       uint64
	Available   uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
	SwapTotal   uint64
	SwapUsed    uint64
	SwapFree    uint64
	SwapPercent float64
}

type memoryBusiness struct{}

var MemoryBusiness = &memoryBusiness{}

func (b *memoryBusiness) Usage() (*MemoryInfo, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("获取内存信息失败: %v", err)
	}

	swapStat, err := mem.SwapMemory()
	if err != nil {
		return nil, fmt.Errorf("获取交换内存信息失败: %v", err)
	}

	return &MemoryInfo{
		Total:       vmStat.Total,
		Available:   vmStat.Available,
		Used:        vmStat.Used,
		Free:        vmStat.Free,
		UsedPercent: vmStat.UsedPercent,
		SwapTotal:   swapStat.Total,
		SwapUsed:    swapStat.Used,
		SwapFree:    swapStat.Free,
		SwapPercent: swapStat.UsedPercent,
	}, nil
}

func (b *memoryBusiness) FormatMemoryTable(info *MemoryInfo) string {
	var sb strings.Builder

	sb.WriteString("========================================\n")
	sb.WriteString("              内存使用情况              \n")
	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("总内存:       %s\n", b.formatBytes(info.Total)))
	sb.WriteString(fmt.Sprintf("已使用:       %s (%.1f%%)\n", b.formatBytes(info.Used), info.UsedPercent))
	sb.WriteString(fmt.Sprintf("可用:         %s\n", b.formatBytes(info.Available)))
	sb.WriteString(fmt.Sprintf("空闲:         %s\n", b.formatBytes(info.Free)))
	sb.WriteString("========================================\n")
	sb.WriteString("             交换内存情况               \n")
	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("总交换:       %s\n", b.formatBytes(info.SwapTotal)))
	sb.WriteString(fmt.Sprintf("已使用:       %s (%.1f%%)\n", b.formatBytes(info.SwapUsed), info.SwapPercent))
	sb.WriteString(fmt.Sprintf("空闲:         %s\n", b.formatBytes(info.SwapFree)))
	sb.WriteString("========================================\n")

	return sb.String()
}

func (b *memoryBusiness) formatBytes(bytes uint64) string {
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

func (b *memoryBusiness) GetUsedPercent() (float64, error) {
	info, err := b.Usage()
	if err != nil {
		return 0, err
	}
	return info.UsedPercent, nil
}
