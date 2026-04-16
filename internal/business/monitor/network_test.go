package monitor

import (
	"testing"
)

func TestNetworkBusiness_Connections(t *testing.T) {
	connections, err := NetworkBusiness.Connections()
	if err != nil {
		t.Fatalf("Connections failed: %v", err)
	}

	t.Logf("Found %d connections", len(connections))
}

func TestNetworkBusiness_Ports(t *testing.T) {
	ports, err := NetworkBusiness.Ports()
	if err != nil {
		t.Fatalf("Ports failed: %v", err)
	}

	t.Logf("Found %d ports", len(ports))

	for i, port := range ports {
		if i >= 5 {
			break
		}
		t.Logf("Port %d: %d/%s, PID: %d, Process: %s",
			i, port.Port, port.Proto, port.Pid, port.Process)
	}
}

func TestNetworkBusiness_FormatConnectionTable(t *testing.T) {
	connections := []ConnectionInfo{
		{
			Pid:    1234,
			Type:   "TCP",
			Laddr:  "127.0.0.1:8080",
			Raddr:  "0.0.0.0:0",
			Status: "LISTEN",
		},
		{
			Pid:    5678,
			Type:   "TCP",
			Laddr:  "192.168.1.1:443",
			Raddr:  "10.0.0.1:54321",
			Status: "ESTABLISHED",
		},
	}

	output := NetworkBusiness.FormatConnectionTable(connections)
	if output == "" {
		t.Error("FormatConnectionTable returned empty string")
	}

	if len(output) < 50 {
		t.Errorf("FormatConnectionTable output too short: %d", len(output))
	}
}

func TestNetworkBusiness_FormatConnectionTableEmpty(t *testing.T) {
	connections := []ConnectionInfo{}
	output := NetworkBusiness.FormatConnectionTable(connections)

	if output != "没有找到网络连接" {
		t.Errorf("Expected '没有找到网络连接', got '%s'", output)
	}
}

func TestNetworkBusiness_FormatPortTable(t *testing.T) {
	ports := []PortInfo{
		{
			Port:    80,
			Proto:   "TCP",
			State:   "LISTEN",
			Pid:     1234,
			Process: "nginx",
		},
		{
			Port:    443,
			Proto:   "TCP",
			State:   "LISTEN",
			Pid:     1234,
			Process: "nginx",
		},
	}

	output := NetworkBusiness.FormatPortTable(ports)
	if output == "" {
		t.Error("FormatPortTable returned empty string")
	}

	if len(output) < 50 {
		t.Errorf("FormatPortTable output too short: %d", len(output))
	}
}

func TestNetworkBusiness_FormatPortTableEmpty(t *testing.T) {
	ports := []PortInfo{}
	output := NetworkBusiness.FormatPortTable(ports)

	if output != "没有找到监听端口" {
		t.Errorf("Expected '没有找到监听端口', got '%s'", output)
	}
}

func TestNetworkBusiness_GetFamilyString(t *testing.T) {
	tests := []struct {
		family   uint32
		expected string
	}{
		{2, "IPv4"},
		{10, "IPv6"},
		{99, "Unknown"},
	}

	for _, test := range tests {
		result := NetworkBusiness.getFamilyString(test.family)
		if result != test.expected {
			t.Errorf("getFamilyString(%d) = %s, expected %s", test.family, result, test.expected)
		}
	}
}

func TestNetworkBusiness_GetTypeString(t *testing.T) {
	tests := []struct {
		connType uint32
		expected string
	}{
		{1, "TCP"},
		{2, "UDP"},
		{99, "Unknown"},
	}

	for _, test := range tests {
		result := NetworkBusiness.getTypeString(test.connType)
		if result != test.expected {
			t.Errorf("getTypeString(%d) = %s, expected %s", test.connType, result, test.expected)
		}
	}
}
