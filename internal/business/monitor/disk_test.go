package monitor

import (
	"testing"
)

func TestDiskBusiness_Usage(t *testing.T) {
	disks, err := DiskBusiness.Usage()
	if err != nil {
		t.Fatalf("Usage failed: %v", err)
	}

	if len(disks) == 0 {
		t.Error("Expected at least one disk partition")
	}

	validDiskCount := 0
	for _, disk := range disks {
		if disk.Total > 0 {
			validDiskCount++
		}
	}

	if validDiskCount == 0 {
		t.Error("Expected at least one disk with non-zero total size")
	}

	t.Logf("Found %d disk partitions, %d with valid size", len(disks), validDiskCount)
}

func TestDiskBusiness_List(t *testing.T) {
	disks, err := DiskBusiness.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(disks) == 0 {
		t.Error("Expected at least one disk partition")
	}
}

func TestDiskBusiness_FormatDiskTable(t *testing.T) {
	disks := []DiskInfo{
		{
			Device:      "/dev/disk1",
			MountPoint:  "/",
			Total:       1073741824,
			Used:        536870912,
			Free:        536870912,
			UsedPercent: 50.0,
			Fstype:      "apfs",
		},
	}

	output := DiskBusiness.FormatDiskTable(disks)
	if output == "" {
		t.Error("FormatDiskTable returned empty string")
	}

	if len(output) < 50 {
		t.Errorf("FormatDiskTable output too short: %d", len(output))
	}
}

func TestDiskBusiness_FormatDiskTableEmpty(t *testing.T) {
	disks := []DiskInfo{}
	output := DiskBusiness.FormatDiskTable(disks)

	if output != "没有找到磁盘分区" {
		t.Errorf("Expected '没有找到磁盘分区', got '%s'", output)
	}
}

func TestDiskBusiness_FormatBytes(t *testing.T) {
	tests := []struct {
		bytes    uint64
		expected string
	}{
		{500, "500B"},
		{1024, "1.00KB"},
		{1048576, "1.00MB"},
		{1073741824, "1.00GB"},
		{1099511627776, "1.00TB"},
	}

	for _, test := range tests {
		result := DiskBusiness.formatBytes(test.bytes)
		if result != test.expected {
			t.Errorf("formatBytes(%d) = %s, expected %s", test.bytes, result, test.expected)
		}
	}
}

func TestDiskBusiness_GetRootDiskUsage(t *testing.T) {
	disk, err := DiskBusiness.GetRootDiskUsage()
	if err != nil {
		t.Fatalf("GetRootDiskUsage failed: %v", err)
	}

	if disk == nil {
		t.Error("Expected disk info, got nil")
		return
	}

	if disk.Total == 0 {
		t.Error("Root disk has zero total size")
	}

	t.Logf("Root disk: %s, Total: %d bytes, Used: %.1f%%",
		disk.MountPoint, disk.Total, disk.UsedPercent)
}
