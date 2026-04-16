package monitor

import (
	"context"

	"github.com/shirou/gopsutil/v3/disk"
)

type DiskUsageInfo struct {
	Path        string  `json:"path"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

type PartitionInfo struct {
	Device     string   `json:"device"`
	Mountpoint string   `json:"mountpoint"`
	Fstype     string   `json:"fstype"`
	Opts       []string `json:"opts"`
}

type DiskBusiness struct{}

var MonitorDisk = &DiskBusiness{}

func (b *DiskBusiness) Usage(ctx context.Context, path string) (*DiskUsageInfo, error) {
	if path == "" {
		path = "/"
	}

	usage, err := disk.Usage(path)
	if err != nil {
		return nil, err
	}

	return &DiskUsageInfo{
		Path:        usage.Path,
		Fstype:      usage.Fstype,
		Total:       usage.Total,
		Free:        usage.Free,
		Used:        usage.Used,
		UsedPercent: usage.UsedPercent,
	}, nil
}

func (b *DiskBusiness) List(ctx context.Context) ([]PartitionInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var result []PartitionInfo
	for _, p := range partitions {
		result = append(result, PartitionInfo{
			Device:     p.Device,
			Mountpoint: p.Mountpoint,
			Fstype:     p.Fstype,
			Opts:       p.Opts,
		})
	}

	return result, nil
}
