package monitor

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/net"
)

// ConnectionInfo 网络连接信息结构体
type ConnectionInfo struct {
	Fd         uint32 // 文件描述符
	Family     string // 协议族(IPv4/IPv6)
	Type       string // 连接类型(TCP/UDP)
	LocalAddr  string // 本地地址
	LocalPort  uint32 // 本地端口
	RemoteAddr string // 远程地址
	RemotePort uint32 // 远程端口
	Status     string // 连接状态
	Pid        int32  // 进程ID
}

// PortInfo 端口信息结构体
type PortInfo struct {
	Port     uint32 // 端口号
	Protocol string // 协议类型
	Status   string // 状态
	Pid      int32  // 进程ID
	Process  string // 进程名称
}

// networkBusiness 网络管理业务逻辑
type networkBusiness struct{}

// NetworkBusiness 网络管理业务实例
var NetworkBusiness = &networkBusiness{}

// Connections 获取网络连接列表
func (b *networkBusiness) Connections() ([]ConnectionInfo, error) {
	connections, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("获取网络连接失败: %v", err)
	}

	var connInfos []ConnectionInfo
	for _, conn := range connections {
		connInfos = append(connInfos, ConnectionInfo{
			Fd:         conn.Fd,
			Family:     b.getFamilyString(conn.Family),
			Type:       b.getTypeString(conn.Type),
			LocalAddr:  conn.Laddr.IP,
			LocalPort:  conn.Laddr.Port,
			RemoteAddr: conn.Raddr.IP,
			RemotePort: conn.Raddr.Port,
			Status:     conn.Status,
			Pid:        conn.Pid,
		})
	}

	return connInfos, nil
}

// Ports 获取监听端口列表
func (b *networkBusiness) Ports() ([]PortInfo, error) {
	connections, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("获取端口列表失败: %v", err)
	}

	portMap := make(map[string]PortInfo)
	for _, conn := range connections {
		if conn.Status == "LISTEN" || conn.Laddr.Port != 0 {
			key := fmt.Sprintf("%d-%s", conn.Laddr.Port, b.getTypeString(conn.Type))
			if _, exists := portMap[key]; !exists {
				portMap[key] = PortInfo{
					Port:     conn.Laddr.Port,
					Protocol: b.getTypeString(conn.Type),
					Status:   conn.Status,
					Pid:      conn.Pid,
				}
			}
		}
	}

	var ports []PortInfo
	for _, port := range portMap {
		ports = append(ports, port)
	}

	return ports, nil
}

// KillPort 杀死占用指定端口的进程
// port: 端口号
func (b *networkBusiness) KillPort(port int) error {
	connections, err := net.Connections("all")
	if err != nil {
		return fmt.Errorf("获取端口列表失败: %v", err)
	}

	var found bool
	for _, conn := range connections {
		if int(conn.Laddr.Port) == port {
			found = true
			if conn.Pid > 0 {
				if runtime.GOOS == "windows" {
					cmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(int(conn.Pid)))
					return cmd.Run()
				}
				cmd := exec.Command("kill", "-9", strconv.Itoa(int(conn.Pid)))
				err := cmd.Run()
				if err != nil {
					return fmt.Errorf("杀死进程失败: %v", err)
				}
				fmt.Printf("✓ 已杀死占用端口 %d 的进程 [PID: %d]\n", port, conn.Pid)
				return nil
			}
		}
	}

	if !found {
		return fmt.Errorf("端口 %d 未被占用", port)
	}

	return fmt.Errorf("无法找到占用端口 %d 的进程", port)
}

// FormatConnectionTable 格式化连接信息为表格字符串
func (b *networkBusiness) FormatConnectionTable(connections []ConnectionInfo) string {
	if len(connections) == 0 {
		return "没有找到网络连接"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-8s %-6s %-20s %-20s %-12s\n", "PID", "类型", "本地地址", "远程地址", "状态"))
	sb.WriteString(strings.Repeat("-", 75) + "\n")

	for _, c := range connections {
		local := fmt.Sprintf("%s:%d", c.LocalAddr, c.LocalPort)
		remote := fmt.Sprintf("%s:%d", c.RemoteAddr, c.RemotePort)
		if c.RemoteAddr == "" {
			remote = "*:*"
		}
		sb.WriteString(fmt.Sprintf("%-8d %-6s %-20s %-20s %-12s\n",
			c.Pid, c.Type, local, remote, c.Status))
	}

	return sb.String()
}

// FormatPortTable 格式化端口信息为表格字符串
func (b *networkBusiness) FormatPortTable(ports []PortInfo) string {
	if len(ports) == 0 {
		return "没有找到监听端口"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-10s %-10s %-12s %-10s\n", "端口", "协议", "状态", "PID"))
	sb.WriteString(strings.Repeat("-", 50) + "\n")

	for _, p := range ports {
		sb.WriteString(fmt.Sprintf("%-10d %-10s %-12s %-10d\n",
			p.Port, p.Protocol, p.Status, p.Pid))
	}

	return sb.String()
}

// getFamilyString 将协议族数值转换为字符串
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

// getTypeString 将连接类型数值转换为字符串
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

// GetProtoString 将协议数值转换为字符串
func (b *networkBusiness) GetProtoString(proto uint32) string {
	switch proto {
	case 6:
		return "TCP"
	case 17:
		return "UDP"
	default:
		return "Unknown"
	}
}

// GetIOCounters 获取网络IO统计信息
func (b *networkBusiness) GetIOCounters() ([]net.IOCountersStat, error) {
	ioCounters, err := net.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("获取网络IO统计失败: %v", err)
	}
	return ioCounters, nil
}
