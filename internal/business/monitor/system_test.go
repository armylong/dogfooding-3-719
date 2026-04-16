package monitor

import (
	"testing"
	"time"
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

	if uptime == nil {
		t.Fatal("Expected uptime info, got nil")
	}

	if uptime.Total == 0 {
		t.Error("Total uptime should not be zero")
	}

	t.Logf("Uptime: %d days, %d hours, %d minutes, %d seconds",
		uptime.Days, uptime.Hours, uptime.Minutes, uptime.Seconds)
}

func TestSystemBusiness_Load(t *testing.T) {
	load, err := SystemBusiness.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if load == nil {
		t.Fatal("Expected load info, got nil")
	}

	t.Logf("Load: 1min=%.2f, 5min=%.2f, 15min=%.2f",
		load.Load1, load.Load5, load.Load15)
}

func TestSystemBusiness_FormatSystemInfo(t *testing.T) {
	info := &SystemInfo{
		Hostname:        "test-host",
		OS:              "darwin",
		Platform:        "macOS",
		PlatformVersion: "14.0",
		KernelVersion:   "23.0.0",
		Architecture:    "arm64",
		Procs:           200,
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
	uptime := &UptimeInfo{
		Days:    1,
		Hours:   2,
		Minutes: 30,
		Seconds: 45,
		Total:   95445,
	}

	output := SystemBusiness.FormatUptime(uptime)
	if output == "" {
		t.Error("FormatUptime returned empty string")
	}

	if len(output) < 50 {
		t.Errorf("FormatUptime output too short: %d", len(output))
	}

	t.Logf("Output:\n%s", output)
}

func TestSystemBusiness_FormatLoad(t *testing.T) {
	load := &LoadInfo{
		Load1:  1.5,
		Load5:  1.2,
		Load15: 1.0,
	}

	output := SystemBusiness.FormatLoad(load)
	if output == "" {
		t.Error("FormatLoad returned empty string")
	}

	if len(output) < 50 {
		t.Errorf("FormatLoad output too short: %d", len(output))
	}

	t.Logf("Output:\n%s", output)
}

func TestSystemBusiness_GetBootTime(t *testing.T) {
	bootTime, err := SystemBusiness.GetBootTime()
	if err != nil {
		t.Fatalf("GetBootTime failed: %v", err)
	}

	if bootTime.IsZero() {
		t.Error("Boot time should not be zero")
	}

	if bootTime.After(time.Now()) {
		t.Error("Boot time should not be in the future")
	}

	t.Logf("Boot time: %s", bootTime.Format("2006-01-02 15:04:05"))
}

func TestSystemBusiness_GetUsers(t *testing.T) {
	users, err := SystemBusiness.GetUsers()
	if err != nil {
		t.Logf("GetUsers failed (may be expected on some systems): %v", err)
		return
	}

	t.Logf("Found %d users", len(users))
}
