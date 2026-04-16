package monitor

import (
	"context"
	"testing"
)

func TestProcessList(t *testing.T) {
	ctx := context.Background()
	processes, err := MonitorProcess.List(ctx)
	if err != nil {
		t.Fatalf("获取进程列表失败: %v", err)
	}

	if len(processes) == 0 {
		t.Error("进程列表不应为空")
	}

	t.Logf("获取到 %d 个进程", len(processes))
}

func TestProcessTop(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name   string
		sortBy string
		limit  int
	}{
		{"按CPU排序", "cpu", 5},
		{"按内存排序", "memory", 5},
		{"按PID排序", "pid", 5},
		{"默认排序", "", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processes, err := MonitorProcess.Top(ctx, tt.sortBy, tt.limit)
			if err != nil {
				t.Fatalf("获取TOP进程失败: %v", err)
			}

			if len(processes) > tt.limit {
				t.Errorf("进程数量超过限制: 期望<=%d, 实际=%d", tt.limit, len(processes))
			}
		})
	}
}

func TestProcessFind(t *testing.T) {
	ctx := context.Background()

	processes, err := MonitorProcess.Find(ctx, "go")
	if err != nil {
		t.Fatalf("查找进程失败: %v", err)
	}

	t.Logf("找到 %d 个包含 'go' 的进程", len(processes))
}

func TestParsePID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int32
		wantErr bool
	}{
		{"有效PID", "1234", 1234, false},
		{"PID为0", "0", 0, false},
		{"无效PID", "abc", 0, true},
		{"负数PID", "-1", -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParsePID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiskUsage(t *testing.T) {
	ctx := context.Background()
	usage, err := MonitorDisk.Usage(ctx, "/")
	if err != nil {
		t.Fatalf("获取磁盘使用情况失败: %v", err)
	}

	if usage.Total == 0 {
		t.Error("磁盘总容量不应为0")
	}

	t.Logf("磁盘总容量: %d, 使用率: %.2f%%", usage.Total, usage.UsedPercent)
}

func TestDiskList(t *testing.T) {
	ctx := context.Background()
	partitions, err := MonitorDisk.List(ctx)
	if err != nil {
		t.Fatalf("获取磁盘分区失败: %v", err)
	}

	if len(partitions) == 0 {
		t.Error("磁盘分区列表不应为空")
	}

	for _, p := range partitions {
		t.Logf("分区: %s, 挂载点: %s", p.Device, p.Mountpoint)
	}
}

func TestMemoryUsage(t *testing.T) {
	ctx := context.Background()
	mem, err := MonitorMemory.Usage(ctx)
	if err != nil {
		t.Fatalf("获取内存使用情况失败: %v", err)
	}

	if mem.Total == 0 {
		t.Error("内存总容量不应为0")
	}

	t.Logf("内存总容量: %d, 使用率: %.2f%%", mem.Total, mem.UsedPercent)
}

func TestCPUUsage(t *testing.T) {
	ctx := context.Background()

	usage, err := MonitorCPU.Usage(ctx, false)
	if err != nil {
		t.Fatalf("获取CPU使用率失败: %v", err)
	}

	t.Logf("CPU平均使用率: %.2f%%", usage.Average)
}

func TestCPUInfo(t *testing.T) {
	ctx := context.Background()
	infos, err := MonitorCPU.Info(ctx)
	if err != nil {
		t.Fatalf("获取CPU信息失败: %v", err)
	}

	if len(infos) == 0 {
		t.Error("CPU信息列表不应为空")
	}

	for _, info := range infos {
		t.Logf("CPU %d: %s", info.CPU, info.ModelName)
	}
}

func TestNetworkConnections(t *testing.T) {
	ctx := context.Background()
	conns, err := MonitorNetwork.Connections(ctx)
	if err != nil {
		t.Fatalf("获取网络连接失败: %v", err)
	}

	t.Logf("获取到 %d 个网络连接", len(conns))
}

func TestNetworkPorts(t *testing.T) {
	ctx := context.Background()
	ports, err := MonitorNetwork.Ports(ctx)
	if err != nil {
		t.Fatalf("获取端口占用失败: %v", err)
	}

	t.Logf("获取到 %d 个占用端口", len(ports))
}

func TestParsePort(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"有效端口", "8080", 8080, false},
		{"端口0", "0", 0, false},
		{"无效端口", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePort(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParsePort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemInfo(t *testing.T) {
	ctx := context.Background()
	info, err := MonitorSystem.Info(ctx)
	if err != nil {
		t.Fatalf("获取系统信息失败: %v", err)
	}

	if info.Hostname == "" {
		t.Error("主机名不应为空")
	}

	t.Logf("主机名: %s, 系统: %s", info.Hostname, info.OS)
}

func TestSystemUptime(t *testing.T) {
	ctx := context.Background()
	uptime, err := MonitorSystem.Uptime(ctx)
	if err != nil {
		t.Fatalf("获取系统运行时间失败: %v", err)
	}

	if uptime.Uptime == 0 {
		t.Error("系统运行时间不应为0")
	}

	t.Logf("系统运行时间: %d天 %d小时 %d分钟", uptime.Days, uptime.Hours, uptime.Minutes)
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name  string
		bytes uint64
	}{
		{"Bytes", 500},
		{"KB", 2 * 1024},
		{"MB", 5 * 1024 * 1024},
		{"GB", 3 * 1024 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			t.Logf("%d -> %s", tt.bytes, result)
		})
	}
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return string(rune(bytes)) + " B"
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return string(rune(exp)) + " B"
}
