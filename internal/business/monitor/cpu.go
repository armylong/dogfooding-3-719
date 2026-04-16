package monitor

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUUsageInfo struct {
	Percent []float64 `json:"percent"`
	Average float64   `json:"average"`
}

type CPUInfo struct {
	CPU        int32   `json:"cpu"`
	VendorID   string  `json:"vendor_id"`
	Family     string  `json:"family"`
	Model      string  `json:"model"`
	ModelName  string  `json:"model_name"`
	Stepping   int32   `json:"stepping"`
	Mhz        float64 `json:"mhz"`
	CacheSize  int32   `json:"cache_size"`
	Cores      int32   `json:"cores"`
	PhysicalID string  `json:"physical_id"`
}

type CPUBusiness struct{}

var MonitorCPU = &CPUBusiness{}

func (b *CPUBusiness) Usage(ctx context.Context, perCPU bool) (*CPUUsageInfo, error) {
	percent, err := cpu.Percent(time.Second, perCPU)
	if err != nil {
		return nil, err
	}

	avg := 0.0
	if len(percent) > 0 {
		for _, p := range percent {
			avg += p
		}
		avg = avg / float64(len(percent))
	}

	return &CPUUsageInfo{
		Percent: percent,
		Average: avg,
	}, nil
}

func (b *CPUBusiness) Info(ctx context.Context) ([]CPUInfo, error) {
	infos, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	var result []CPUInfo
	for _, info := range infos {
		result = append(result, CPUInfo{
			CPU:        info.CPU,
			VendorID:   info.VendorID,
			Family:     info.Family,
			Model:      info.Model,
			ModelName:  info.ModelName,
			Stepping:   info.Stepping,
			Mhz:        info.Mhz,
			CacheSize:  info.CacheSize,
			Cores:      info.Cores,
			PhysicalID: info.PhysicalID,
		})
	}

	return result, nil
}
