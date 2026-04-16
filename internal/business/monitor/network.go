package monitor

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/net"
)

type ConnectionInfo struct {
	Fd     uint32
	Family string
	Type   string
	Laddr  string
	Raddr  string
	Status string
	Pid    int32
}

type PortInfo struct {
	Port    uint32
	Proto   string
	State   string
	Pid     int32
	Process string
}

type networkBusiness struct{}

var NetworkBusiness = &networkBusiness{}

func (b *networkBusiness) Connections() ([]ConnectionInfo, error) {
	connections, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("获取网络连接失败: %v", err)
	}

	var connInfos []ConnectionInfo
	for _, conn := range connections {
		connInfos = append(connInfos, ConnectionInfo{
			Fd:     conn.Fd,
			Family: b.getFamilyString(conn.Family),
			Type:   b.getTypeString(conn.Type),
			Laddr:  fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port),
			Raddr:  fmt.Sprintf("%s:%d", conn.Raddr.IP, conn.Raddr.Port),
			Status: conn.Status,
			Pid:    conn.Pid,
		})
	}

	return connInfos, nil
}

func (b *networkBusiness) Ports() ([]PortInfo, error) {
	connections, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("获取端口信息失败: %v", err)
	}

	portMap := make(map[string]PortInfo)
	for _, conn := range connections {
		if conn.Laddr.Port == 0 {
			continue
		}

		key := fmt.Sprintf("%d-%s", conn.Laddr.Port, b.getProtoString(conn.Type))
		if _, exists := portMap[key]; !exists {
			processName := ""
			if conn.Pid > 0 {
				processName = b.getProcessName(conn.Pid)
			}

			portMap[key] = PortInfo{
				Port:    conn.Laddr.Port,
				Proto:   b.getProtoString(conn.Type),
				State:   conn.Status,
				Pid:     conn.Pid,
				Process: processName,
			}
		}
	}

	var ports []PortInfo
	for _, port := range portMap {
		ports = append(ports, port)
	}

	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Port < ports[j].Port
	})

	return ports, nil
}

func (b *networkBusiness) KillPort(port int) error {
	connections, err := net.Connections("all")
	if err != nil {
		return fmt.Errorf("获取端口信息失败: %v", err)
	}

	var pids []int32
	for _, conn := range connections {
		if conn.Laddr.Port == uint32(port) && conn.Pid > 0 {
			pids = append(pids, conn.Pid)
		}
	}

	if len(pids) == 0 {
		return fmt.Errorf("没有找到占用端口 %d 的进程", port)
	}

	for _, pid := range pids {
		if runtime.GOOS == "windows" {
			cmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(int(pid)))
			if err := cmd.Run(); err != nil {
				fmt.Printf("杀死进程 %d 失败: %v\n", pid, err)
			} else {
				fmt.Printf("✓ 已杀死进程 [PID: %d] 占用端口 %d\n", pid, port)
			}
		} else {
			cmd := exec.Command("kill", "-9", strconv.Itoa(int(pid)))
			if err := cmd.Run(); err != nil {
				fmt.Printf("杀死进程 %d 失败: %v\n", pid, err)
			} else {
				fmt.Printf("✓ 已杀死进程 [PID: %d] 占用端口 %d\n", pid, port)
			}
		}
	}

	return nil
}

func (b *networkBusiness) FormatConnectionTable(connections []ConnectionInfo) string {
	if len(connections) == 0 {
		return "没有找到网络连接"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-8s %-10s %-25s %-25s %-15s\n", "PID", "协议", "本地地址", "远程地址", "状态"))
	sb.WriteString(strings.Repeat("-", 90) + "\n")

	for _, conn := range connections {
		if conn.Pid == 0 && conn.Status == "" {
			continue
		}
		sb.WriteString(fmt.Sprintf("%-8d %-10s %-25s %-25s %-15s\n",
			conn.Pid, conn.Type, conn.Laddr, conn.Raddr, conn.Status))
	}

	return sb.String()
}

func (b *networkBusiness) FormatPortTable(ports []PortInfo) string {
	if len(ports) == 0 {
		return "没有找到监听端口"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-8s %-10s %-15s %-10s %-20s\n", "端口", "协议", "状态", "PID", "进程"))
	sb.WriteString(strings.Repeat("-", 75) + "\n")

	for _, port := range ports {
		sb.WriteString(fmt.Sprintf("%-8d %-10s %-15s %-10d %-20s\n",
			port.Port, port.Proto, port.State, port.Pid, port.Process))
	}

	return sb.String()
}

func (b *networkBusiness) getFamilyString(family uint32) string {
	switch family {
	case 2:
		return "IPv4"
	case 10:
		return "IPv6"
	default:
		return "Unknown"
	}
}

func (b *networkBusiness) getTypeString(connType uint32) string {
	switch connType {
	case 1:
		return "TCP"
	case 2:
		return "UDP"
	default:
		return "Unknown"
	}
}

func (b *networkBusiness) getProtoString(connType uint32) string {
	switch connType {
	case 1:
		return "TCP"
	case 2:
		return "UDP"
	default:
		return "Unknown"
	}
}

func (b *networkBusiness) getProcessName(pid int32) string {
	procs, err := net.Connections("all")
	if err != nil {
		return ""
	}

	for _, proc := range procs {
		if proc.Pid == pid {
			return fmt.Sprintf("PID:%d", pid)
		}
	}
	return ""
}
