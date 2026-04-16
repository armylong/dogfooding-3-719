package monitor

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUInfo struct {
	ModelName string
	Cores     int
	Usage     float64
}

type CPUUsage struct {
	User    float64
	System  float64
	Idle    float64
	Nice    float64
	Iowait  float64
	Irq     float64
	Softirq float64
}

type cpuBusiness struct{}

var CPUBusiness = &cpuBusiness{}

func (b *cpuBusiness) Usage(interval time.Duration) ([]float64, error) {
	percentages, err := cpu.Percent(interval, true)
	if err != nil {
		return nil, fmt.Errorf("获取CPU使用率失败: %v", err)
	}
	return percentages, nil
}

func (b *cpuBusiness) Info() ([]cpu.InfoStat, error) {
	info, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("获取CPU信息失败: %v", err)
	}
	return info, nil
}

func (b *cpuBusiness) Count(logical bool) (int, error) {
	count, err := cpu.Counts(logical)
	if err != nil {
		return 0, fmt.Errorf("获取CPU核心数失败: %v", err)
	}
	return count, nil
}

func (b *cpuBusiness) FormatCPUTable(usage []float64, info []cpu.InfoStat) string {
	var sb strings.Builder

	sb.WriteString("========================================\n")
	sb.WriteString("              CPU 信息                  \n")
	sb.WriteString("========================================\n")

	if len(info) > 0 {
		sb.WriteString(fmt.Sprintf("型号:         %s\n", info[0].ModelName))
		sb.WriteString(fmt.Sprintf("核心数:       %d\n", len(info)))
		sb.WriteString(fmt.Sprintf("物理核心:     %d\n", info[0].Cores))
	}

	logicalCores, _ := b.Count(true)
	sb.WriteString(fmt.Sprintf("逻辑核心:     %d\n", logicalCores))
	sb.WriteString("========================================\n")
	sb.WriteString("            各核心使用率                \n")
	sb.WriteString("========================================\n")

	for i, u := range usage {
		sb.WriteString(fmt.Sprintf("核心 %d:       %.1f%%\n", i, u))
	}

	avgUsage := b.calculateAverage(usage)
	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("平均使用率:   %.1f%%\n", avgUsage))
	sb.WriteString("========================================\n")

	return sb.String()
}

func (b *cpuBusiness) FormatUsageOnly(usage []float64) string {
	var sb strings.Builder

	sb.WriteString("========================================\n")
	sb.WriteString("            CPU 使用率                  \n")
	sb.WriteString("========================================\n")

	for i, u := range usage {
		sb.WriteString(fmt.Sprintf("核心 %d:       %.1f%%\n", i, u))
	}

	avgUsage := b.calculateAverage(usage)
	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("平均使用率:   %.1f%%\n", avgUsage))
	sb.WriteString("========================================\n")

	return sb.String()
}

func (b *cpuBusiness) calculateAverage(usage []float64) float64 {
	if len(usage) == 0 {
		return 0
	}
	var total float64
	for _, u := range usage {
		total += u
	}
	return total / float64(len(usage))
}

func (b *cpuBusiness) GetTotalUsage(interval time.Duration) (float64, error) {
	percentages, err := cpu.Percent(interval, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) > 0 {
		return percentages[0], nil
	}
	return 0, nil
}
