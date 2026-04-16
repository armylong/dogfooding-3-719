package monitor

import (
	"context"
	"fmt"
	gonet "net"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type ConnectionInfo struct {
	Fd        uint32 `json:"fd"`
	Family    uint32 `json:"family"`
	Type      uint32 `json:"type"`
	Laddr     string `json:"laddr"`
	Raddr     string `json:"raddr"`
	Status    string `json:"status"`
	PID       int32  `json:"pid"`
	LocalIP   string `json:"local_ip"`
	LocalPort uint32 `json:"local_port"`
	RemoteIP  string `json:"remote_ip"`
	RemotePort uint32 `json:"remote_port"`
}

type PortInfo struct {
	Port     uint32 `json:"port"`
	Protocol string `json:"protocol"`
	Status   string `json:"status"`
	PID      int32  `json:"pid"`
	Name     string `json:"name"`
}

type NetworkBusiness struct{}

var MonitorNetwork = &NetworkBusiness{}

func (b *NetworkBusiness) Connections(ctx context.Context) ([]ConnectionInfo, error) {
	conns, err := net.Connections("all")
	if err != nil {
		return nil, err
	}

	var result []ConnectionInfo
	for _, c := range conns {
		info := ConnectionInfo{
			Fd:     c.Fd,
			Family: c.Family,
			Type:   c.Type,
			Status: c.Status,
			PID:    c.Pid,
		}

		if c.Laddr.Port > 0 || c.Laddr.IP != "" {
			info.Laddr = fmt.Sprintf("%s:%d", c.Laddr.IP, c.Laddr.Port)
			info.LocalIP = c.Laddr.IP
			info.LocalPort = c.Laddr.Port
		}

		if c.Raddr.Port > 0 || c.Raddr.IP != "" {
			info.Raddr = fmt.Sprintf("%s:%d", c.Raddr.IP, c.Raddr.Port)
			info.RemoteIP = c.Raddr.IP
			info.RemotePort = c.Raddr.Port
		}

		result = append(result, info)
	}

	return result, nil
}

func (b *NetworkBusiness) Ports(ctx context.Context) ([]PortInfo, error) {
	conns, err := net.Connections("all")
	if err != nil {
		return nil, err
	}

	portMap := make(map[string]PortInfo)
	var result []PortInfo

	for _, c := range conns {
		if c.Laddr.Port > 0 {
			key := fmt.Sprintf("%d-%d", c.Laddr.Port, c.Type)
			if _, exists := portMap[key]; !exists {
				proto := "TCP"
				if c.Type == 2 {
					proto = "UDP"
				}

				name := ""
				if c.Pid > 0 {
					name = getProcessName(c.Pid)
				}

				portMap[key] = PortInfo{
					Port:     c.Laddr.Port,
					Protocol: proto,
					Status:   c.Status,
					PID:      c.Pid,
					Name:     name,
				}
			}
		}
	}

	for _, p := range portMap {
		result = append(result, p)
	}

	return result, nil
}

func (b *NetworkBusiness) KillPort(ctx context.Context, port int) error {
	conns, err := net.Connections("all")
	if err != nil {
		return err
	}

	pidMap := make(map[int32]bool)
	for _, c := range conns {
		if c.Laddr.Port == uint32(port) && c.Pid > 0 {
			pidMap[c.Pid] = true
		}
	}

	if len(pidMap) == 0 {
		return fmt.Errorf("no process found on port %d", port)
	}

	for pid := range pidMap {
		p, err := process.NewProcess(pid)
		if err != nil {
			continue
		}
		_ = p.Kill()
	}

	return nil
}

func getProcessName(pid int32) string {
	p, err := process.NewProcess(pid)
	if err != nil {
		return ""
	}
	name, _ := p.Name()
	return name
}

func GetPortProcessPID(port int) (int, error) {
	addr, err := gonet.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return 0, err
	}

	listener, err := gonet.ListenTCP("tcp", addr)
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			conns, _ := net.Connections("tcp")
			for _, c := range conns {
				if c.Laddr.Port == uint32(port) && c.Status == "LISTEN" {
					return int(c.Pid), nil
				}
			}
		}
		return 0, err
	}
	defer listener.Close()

	return 0, nil
}

func ParsePort(portStr string) (int, error) {
	return strconv.Atoi(portStr)
}
