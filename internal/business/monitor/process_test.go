package monitor

import (
	"testing"
)

func TestProcessBusiness_ListProcesses(t *testing.T) {
	processes, err := ProcessBusiness.ListProcesses("cpu", 5)
	if err != nil {
		t.Fatalf("ListProcesses failed: %v", err)
	}

	if len(processes) == 0 {
		t.Error("Expected at least one process")
	}

	if len(processes) > 5 {
		t.Errorf("Expected at most 5 processes, got %d", len(processes))
	}
}

func TestProcessBusiness_SortProcesses(t *testing.T) {
	processes := []ProcessInfo{
		{PID: 1, Name: "proc1", CPU: 10.0, Memory: 20.0},
		{PID: 2, Name: "proc2", CPU: 30.0, Memory: 10.0},
		{PID: 3, Name: "proc3", CPU: 20.0, Memory: 30.0},
	}

	ProcessBusiness.sortProcesses(processes, "cpu")
	if processes[0].CPU != 30.0 {
		t.Errorf("Expected first process CPU to be 30.0, got %f", processes[0].CPU)
	}

	ProcessBusiness.sortProcesses(processes, "memory")
	if processes[0].Memory != 30.0 {
		t.Errorf("Expected first process Memory to be 30.0, got %f", processes[0].Memory)
	}

	ProcessBusiness.sortProcesses(processes, "pid")
	if processes[0].PID != 1 {
		t.Errorf("Expected first process PID to be 1, got %d", processes[0].PID)
	}
}

func TestProcessBusiness_FormatProcessTable(t *testing.T) {
	processes := []ProcessInfo{
		{PID: 1, Name: "test", CPU: 10.0, Memory: 20.0, Status: "running"},
	}

	output := ProcessBusiness.FormatProcessTable(processes)
	if output == "" {
		t.Error("FormatProcessTable returned empty string")
	}

	if len(output) < 50 {
		t.Errorf("FormatProcessTable output too short: %d", len(output))
	}
}

func TestProcessBusiness_FormatProcessTableEmpty(t *testing.T) {
	processes := []ProcessInfo{}
	output := ProcessBusiness.FormatProcessTable(processes)

	if output != "没有找到进程" {
		t.Errorf("Expected '没有找到进程', got '%s'", output)
	}
}

func TestProcessBusiness_FindProcess(t *testing.T) {
	processes, err := ProcessBusiness.FindProcess("kernel")
	if err != nil {
		t.Fatalf("FindProcess failed: %v", err)
	}

	t.Logf("Found %d processes matching 'kernel'", len(processes))
}

func TestProcessBusiness_GetProcessCount(t *testing.T) {
	count, err := ProcessBusiness.GetProcessCount()
	if err != nil {
		t.Fatalf("GetProcessCount failed: %v", err)
	}

	if count <= 0 {
		t.Errorf("Expected positive process count, got %d", count)
	}

	t.Logf("Process count: %d", count)
}
