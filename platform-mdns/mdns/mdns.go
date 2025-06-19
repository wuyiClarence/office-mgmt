package mdns

import (
	"encoding/base64"
	"fmt"
	"net"
	"sync"

	"platform-mdns/config"
	"platform-mdns/dto"
	"platform-mdns/utils"
	mylog "platform-mdns/utils/log"
)

type IfaceInfo struct {
	Name  string // e.g., "en0", "lo0", "eth0.100"
	Ips   []net.IP
	Iface net.Interface
}

func getInterfaceIPs() (map[string]IfaceInfo, error) {
	interfaceIPs := make(map[string]IfaceInfo)

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		// 排除回环网卡和无效网卡
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if iface.Flags&net.FlagMulticast == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		ifaceInfo := IfaceInfo{
			Name:  iface.Name,
			Ips:   make([]net.IP, 0),
			Iface: iface,
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.To4() != nil {
				ifaceInfo.Ips = append(ifaceInfo.Ips, ipNet.IP)
			}
		}
		if len(ifaceInfo.Ips) == 0 {
			continue
		}
		interfaceIPs[iface.Name] = ifaceInfo
	}

	return interfaceIPs, nil
}

// 获取所有有效的本地 IPv4 地址，返回 net.IP 切片
func getLocalIPv4() ([]net.IP, error) {
	var ipList []net.IP

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		// 忽略未启用的或环回接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && ipNet.IP.To4() != nil {
				ipList = append(ipList, ipNet.IP) // 收集所有可用的 IPv4 地址
			}
		}
	}

	if len(ipList) == 0 {
		return nil, fmt.Errorf("未找到有效的 IPv4 地址")
	}

	return ipList, nil
}

func StartMdns() ([]*Server, error) {
	interfaceIPs, err := getInterfaceIPs() // 获取网卡和对应的 IP 地址
	if err != nil {
		return nil, err
	}
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		servers []*Server
		errChan = make(chan error, len(interfaceIPs)) // 存储错误
	)

	domain := config.MyConfig.MdnsConfig.Domain
	host := config.MyConfig.MdnsConfig.Host
	key := config.MyConfig.MdnsConfig.Key

	MdnsInfo := dto.MdnsInfo{
		Broker:   config.MyConfig.MqttConfig.Broker,
		UserName: config.MyConfig.MqttConfig.UserName,
		Password: config.MyConfig.MqttConfig.Password,
		Port:     config.MyConfig.MqttConfig.Port,
	}

	jsonString, err := utils.ToJSONString(MdnsInfo)
	if err != nil {
		return nil, err
	}
	_, _ = fmt.Fprintf(mylog.SystemLogger, "StartMdns ...\n")
	encryptedBytes := utils.XorEncryptDecrypt(jsonString, key)
	encryptedBase64 := base64.StdEncoding.EncodeToString(encryptedBytes)

	info := []string{"base64=" + encryptedBase64}

	// 并发启动 mDNS 服务
	for ifaceName, ifaceinfo := range interfaceIPs {
		wg.Add(1)
		go func(ifaceName string, ifaceinfo IfaceInfo) {
			defer wg.Done()

			// mDNS 主机名
			hostname := fmt.Sprintf("%s-%s.%s", host, ifaceName, domain)

			service, err := NewMDNSService(host, config.MyConfig.MdnsConfig.ServiceName, domain, hostname, config.MyConfig.MdnsConfig.Port, ifaceinfo.Ips, info, true)
			if err != nil {
				_, _ = fmt.Fprintf(mylog.SystemLogger, "mDNS 服务创建失败: 网卡=%s, IP=%v, err=%v\n", ifaceName, ifaceinfo.Ips, err)
				errChan <- err
				return
			}

			// 启动 mDNS 服务器
			server, err := NewServer(&Config{Zone: service, Iface: &ifaceinfo.Iface})
			if err != nil {
				_, _ = fmt.Fprintf(mylog.SystemLogger, "mDNS 服务器启动失败: 网卡=%s, IP=%v, err=%v\n", ifaceName, ifaceinfo.Ips, err)
				errChan <- err
				return
			}

			mu.Lock()
			servers = append(servers, server)
			mu.Unlock()

			_, _ = fmt.Fprintf(mylog.SystemLogger, "✅ mDNS 服务器启动成功: 网卡=%s, IP=%v, Hostname=%s\n", ifaceName, ifaceinfo.Ips, hostname)
		}(ifaceName, ifaceinfo)
	}
	// 等待所有 goroutine 完成
	wg.Wait()
	close(errChan)
	// 检查错误
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}
	return nil, err
}
