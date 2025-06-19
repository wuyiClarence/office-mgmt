package protocol

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"platform-backend/config"
	"platform-backend/utils/encrypt"

	"github.com/hashicorp/mdns"
)

type MdnsInfo struct {
	Broker   string `json:"mqtt_broker"`
	UserName string `json:"mqtt_username"`
	Password string `json:"mqtt_password"`
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

// 将结构体转换为 JSON 字符串
func toJSONString(info MdnsInfo) (string, error) {
	jsonData, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
func StartMdns() (*mdns.Server, error) {
	// 获取本地 IPv4 地址列表
	ipList, err := getLocalIPv4()
	if err != nil {
		log.Fatalf("无法获取本地 IP 地址: %v", err)
	}
	domain := "local."
	host := "officemgmt"
	MdnsInfo := MdnsInfo{
		Broker:   config.MyConfig.MqttConfig.Broker,
		UserName: config.MyConfig.MqttConfig.UserName,
		Password: config.MyConfig.MqttConfig.Password,
	}

	jsonString, err := toJSONString(MdnsInfo)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil, err
	}
	key := "zwxlink"
	encryptedBytes := encrypt.XorEncryptDecrypt(jsonString, key)
	encryptedBase64 := base64.StdEncoding.EncodeToString(encryptedBytes)
	fmt.Println("Encrypted JSON (Base64):", encryptedBase64)

	info := []string{"base64=" + encryptedBase64}
	hostname := host + "." + domain

	service, err := mdns.NewMDNSService(host, "_officemgmt._tcp", domain, hostname, 80, ipList, info)
	if err != nil {
		log.Fatalf("Unable NewMDNSService: %v", err)
		return nil, err
	}
	// Create the mDNS server, defer shutdown
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		log.Fatalf("Unable NewServer: %v", err)
	}
	//defer server.Shutdown()
	return server, err
}
