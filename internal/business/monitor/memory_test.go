package monitor

import (
	"testing"
)

func TestMemoryBusiness_Usage(t *testing.T) {
	info, err := MemoryBusiness.Usage()
	if err != nil {
		t.Fatalf("Usage failed: %v", err)
	}

	if info == nil {
		t.Fatal("Expected memory info, got nil")
	}

	if info.Total == 0 {
		t.Error("Total memory should not be zero")
	}

	t.Logf("Memory: Total=%d, Used=%.1f%%, Available=%d",
		info.Total, info.UsedPercent, info.Available)
}

func TestMemoryBusiness_FormatMemoryTable(t *testing.T) {
	info := &MemoryInfo{
		Total:       17179869184,
		Available:   8589934592,
		Used:        8589934592,
		Free:        4294967296,
		UsedPercent: 50.0,
		SwapTotal:   4294967296,
		SwapUsed:    2147483648,
		SwapFree:    2147483648,
		SwapPercent: 50.0,
	}

	output := MemoryBusiness.FormatMemoryTable(info)
	if output == "" {
		t.Error("FormatMemoryTable returned empty string")
	}

	if len(output) < 100 {
		t.Errorf("FormatMemoryTable output too short: %d", len(output))
	}

	t.Logf("Output:\n%s", output)
}

func TestMemoryBusiness_FormatBytes(t *testing.T) {
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
		result := MemoryBusiness.formatBytes(test.bytes)
		if result != test.expected {
			t.Errorf("formatBytes(%d) = %s, expected %s", test.bytes, result, test.expected)
		}
	}
}

func TestMemoryBusiness_GetUsedPercent(t *testing.T) {
	percent, err := MemoryBusiness.GetUsedPercent()
	if err != nil {
		t.Fatalf("GetUsedPercent failed: %v", err)
	}

	if percent < 0 || percent > 100 {
		t.Errorf("Invalid memory percentage: %f", percent)
	}

	t.Logf("Memory used percent: %.1f%%", percent)
}
