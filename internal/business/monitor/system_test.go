package monitor

import (
	"testing"
)

func TestSystemBusiness_Info(t *testing.T) {
	info, err := SystemBusiness.Info()
	if err != nil {
		t.Fatalf("Info failed: %v", err)
	}

	if info == nil {
		t.Fatal("Expected system info, got nil")
	}

	if info.Hostname == "" {
		t.Error("Hostname should not be empty")
	}

	if info.OS == "" {
		t.Error("OS should not be empty")
	}

	t.Logf("System: %s %s, Hostname: %s", info.OS, info.Platform, info.Hostname)
}

func TestSystemBusiness_Uptime(t *testing.T) {
	uptime, err := SystemBusiness.Uptime()
	if err != nil {
		t.Fatalf("Uptime failed: %v", err)
	}

	if uptime == 0 {
		t.Error("Uptime should not be zero")
	}

	t.Logf("Uptime: %d seconds", uptime)
}

func TestSystemBusiness_FormatSystemInfo(t *testing.T) {
	info := &SystemInfo{
		Hostname:        "test-host",
		OS:              "darwin",
		Platform:        "macOS",
		PlatformVersion: "14.0",
		KernelVersion:   "23.0.0",
		Architecture:    "arm64",
		Uptime:          95445,
	}

	output := SystemBusiness.FormatSystemInfo(info)
	if output == "" {
		t.Error("FormatSystemInfo returned empty string")
	}

	if len(output) < 50 {
		t.Errorf("FormatSystemInfo output too short: %d", len(output))
	}

	t.Logf("Output:\n%s", output)
}

func TestSystemBusiness_FormatUptime(t *testing.T) {
	uptime := uint64(95445)

	output := SystemBusiness.FormatUptime(uptime)
	if output == "" {
		t.Error("FormatUptime returned empty string")
	}

	t.Logf("Output: %s", output)
}

func TestSystemBusiness_BootTime(t *testing.T) {
	bootTime, err := SystemBusiness.BootTime()
	if err != nil {
		t.Fatalf("BootTime failed: %v", err)
	}

	if bootTime == 0 {
		t.Error("Boot time should not be zero")
	}

	t.Logf("Boot time (unix): %d", bootTime)
}

func TestSystemBusiness_GetBootTimeFormatted(t *testing.T) {
	bootTime, err := SystemBusiness.GetBootTimeFormatted()
	if err != nil {
		t.Fatalf("GetBootTimeFormatted failed: %v", err)
	}

	if bootTime == "" {
		t.Error("Boot time should not be empty")
	}

	t.Logf("Boot time: %s", bootTime)
}

func TestSystemBusiness_Users(t *testing.T) {
	users, err := SystemBusiness.Users()
	if err != nil {
		t.Logf("Users failed (may be expected on some systems): %v", err)
		return
	}

	t.Logf("Found %d users", len(users))
}

func TestSystemBusiness_formatUptime(t *testing.T) {
	tests := []struct {
		seconds  uint64
		contains string
	}{
		{0, "0秒"},
		{60, "1分钟"},
		{3600, "1小时"},
		{86400, "1天"},
		{90061, "1天"},
	}

	for _, test := range tests {
		result := SystemBusiness.formatUptime(test.seconds)
		if result == "" {
			t.Errorf("formatUptime(%d) returned empty string", test.seconds)
		}
		t.Logf("formatUptime(%d) = %s", test.seconds, result)
	}
}
