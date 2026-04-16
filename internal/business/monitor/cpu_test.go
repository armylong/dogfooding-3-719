package monitor

import (
	"testing"
	"time"
)

func TestCPUBusiness_Usage(t *testing.T) {
	usage, err := CPUBusiness.Usage(time.Second)
	if err != nil {
		t.Fatalf("Usage failed: %v", err)
	}

	if len(usage) == 0 {
		t.Error("Expected at least one CPU core")
	}

	for i, u := range usage {
		if u < 0 || u > 100 {
			t.Errorf("Invalid CPU usage for core %d: %f", i, u)
		}
	}

	t.Logf("CPU usage: %d cores", len(usage))
}

func TestCPUBusiness_Info(t *testing.T) {
	info, err := CPUBusiness.Info()
	if err != nil {
		t.Fatalf("Info failed: %v", err)
	}

	if len(info) == 0 {
		t.Error("Expected CPU info")
	}

	for i, cpu := range info {
		t.Logf("CPU %d: %s, Cores: %d", i, cpu.ModelName, cpu.Cores)
	}
}

func TestCPUBusiness_Count(t *testing.T) {
	logical, err := CPUBusiness.Count(true)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	if logical <= 0 {
		t.Errorf("Expected positive CPU count, got %d", logical)
	}

	physical, err := CPUBusiness.Count(false)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	if physical <= 0 {
		t.Errorf("Expected positive physical CPU count, got %d", physical)
	}

	t.Logf("Logical cores: %d, Physical cores: %d", logical, physical)
}

func TestCPUBusiness_FormatCPUTable(t *testing.T) {
	usage := []float64{10.5, 20.3, 15.7, 8.2}

	output := CPUBusiness.FormatCPUTable(usage, nil)
	if output == "" {
		t.Error("FormatCPUTable returned empty string")
	}

	t.Logf("Output:\n%s", output)
}

func TestCPUBusiness_FormatUsageOnly(t *testing.T) {
	usage := []float64{10.5, 20.3, 15.7, 8.2}

	output := CPUBusiness.FormatUsageOnly(usage)
	if output == "" {
		t.Error("FormatUsageOnly returned empty string")
	}

	if len(output) < 50 {
		t.Errorf("FormatUsageOnly output too short: %d", len(output))
	}
}

func TestCPUBusiness_CalculateAverage(t *testing.T) {
	tests := []struct {
		usage    []float64
		expected float64
	}{
		{[]float64{10.0, 20.0, 30.0}, 20.0},
		{[]float64{5.0}, 5.0},
		{[]float64{}, 0.0},
		{[]float64{0.0, 0.0, 0.0}, 0.0},
	}

	for _, test := range tests {
		result := CPUBusiness.calculateAverage(test.usage)
		if result != test.expected {
			t.Errorf("calculateAverage(%v) = %f, expected %f", test.usage, result, test.expected)
		}
	}
}

func TestCPUBusiness_GetTotalUsage(t *testing.T) {
	usage, err := CPUBusiness.GetTotalUsage(time.Second)
	if err != nil {
		t.Fatalf("GetTotalUsage failed: %v", err)
	}

	if usage < 0 || usage > 100 {
		t.Errorf("Invalid total CPU usage: %f", usage)
	}

	t.Logf("Total CPU usage: %.1f%%", usage)
}
