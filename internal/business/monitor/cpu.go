package monitor

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

// CPUInfo CPU信息结构体
type CPUInfo struct {
	ModelName string  // CPU型号
	Cores     int32   // 核心数
	Mhz       float64 // 频率(MHz)
}

// cpuBusiness CPU管理业务逻辑
type cpuBusiness struct{}

// CPUBusiness CPU管理业务实例
var CPUBusiness = &cpuBusiness{}

// Usage 获取CPU使用率
// interval: 采样间隔时间
func (b *cpuBusiness) Usage(interval time.Duration) ([]float64, error) {
	percent, err := cpu.Percent(interval, true)
	if err != nil {
		return nil, fmt.Errorf("获取CPU使用率失败: %v", err)
	}
	return percent, nil
}

// Info 获取CPU信息
func (b *cpuBusiness) Info() ([]CPUInfo, error) {
	cpuInfos, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("获取CPU信息失败: %v", err)
	}

	var infos []CPUInfo
	for _, info := range cpuInfos {
		infos = append(infos, CPUInfo{
			ModelName: info.ModelName,
			Cores:     info.Cores,
			Mhz:       info.Mhz,
		})
	}

	return infos, nil
}

// Count 获取CPU核心数
func (b *cpuBusiness) Count() (int, int, error) {
	logical := runtime.NumCPU()
	physical, err := cpu.Counts(false)
	if err != nil {
		return 0, 0, fmt.Errorf("获取CPU核心数失败: %v", err)
	}
	return logical, physical, nil
}

// FormatCPUTable 格式化CPU信息为表格字符串
func (b *cpuBusiness) FormatCPUTable(usage []float64, info []CPUInfo) string {
	var sb strings.Builder

	sb.WriteString("========================================\n")
	sb.WriteString("              CPU 信息                  \n")
	sb.WriteString("========================================\n")

	if len(info) > 0 {
		sb.WriteString(fmt.Sprintf("型号:         %s\n", info[0].ModelName))
		sb.WriteString(fmt.Sprintf("核心数:       %d\n", info[0].Cores))
	}

	logical, physical, _ := b.Count()
	sb.WriteString(fmt.Sprintf("逻辑核心:     %d\n", logical))
	sb.WriteString(fmt.Sprintf("物理核心:     %d\n", physical))

	sb.WriteString("========================================\n")
	sb.WriteString("            各核心使用率                \n")
	sb.WriteString("========================================\n")

	for i, u := range usage {
		sb.WriteString(fmt.Sprintf("核心 %d:       %.1f%%\n", i, u))
	}

	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("平均使用率:   %.1f%%\n", b.calculateAverage(usage)))
	sb.WriteString("========================================\n")

	return sb.String()
}

// FormatUsageOnly 仅格式化CPU使用率
func (b *cpuBusiness) FormatUsageOnly(usage []float64) string {
	var sb strings.Builder

	sb.WriteString("========================================\n")
	sb.WriteString("            CPU 使用率                  \n")
	sb.WriteString("========================================\n")

	for i, u := range usage {
		sb.WriteString(fmt.Sprintf("核心 %d:       %.1f%%\n", i, u))
	}

	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("平均使用率:   %.1f%%\n", b.calculateAverage(usage)))
	sb.WriteString("========================================\n")

	return sb.String()
}

// calculateAverage 计算平均使用率
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

// GetTotalUsage 获取总体CPU使用率
func (b *cpuBusiness) GetTotalUsage(interval time.Duration) (float64, error) {
	percent, err := cpu.Percent(interval, false)
	if err != nil {
		return 0, fmt.Errorf("获取CPU使用率失败: %v", err)
	}

	if len(percent) > 0 {
		return percent[0], nil
	}

	return 0, nil
}
