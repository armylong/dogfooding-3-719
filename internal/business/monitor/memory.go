package monitor

import (
	"context"

	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryUsageInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Free        uint64  `json:"free"`
	Active      uint64  `json:"active"`
	Inactive    uint64  `json:"inactive"`
	Wired       uint64  `json:"wired"`
	Cached      uint64  `json:"cached"`
}

type MemoryBusiness struct{}

var MonitorMemory = &MemoryBusiness{}

func (b *MemoryBusiness) Usage(ctx context.Context) (*MemoryUsageInfo, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &MemoryUsageInfo{
		Total:       v.Total,
		Available:   v.Available,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
		Free:        v.Free,
		Active:      v.Active,
		Inactive:    v.Inactive,
		Wired:       v.Wired,
		Cached:      v.Cached,
	}, nil
}
