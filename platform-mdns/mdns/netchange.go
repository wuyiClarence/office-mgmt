package mdns

import (
	"fmt"
	"net"
	"os"
	mylog "platform-mdns/utils/log"
	"time"
)

func ipsEqualUnordered(a, b []net.IP) bool {
	if len(a) != len(b) {
		return false
	}
	for _, ipA := range a {
		found := false
		for _, ipB := range b {
			if ipA.Equal(ipB) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
func CheckNetChange() {
	interfaceIPs, err := getInterfaceIPs() // 获取网卡和对应的 IP 地址
	if err != nil {
		_, _ = fmt.Fprintf(mylog.SystemLogger, "CheckNetChange 获取网卡和对应的 IP 地址:%s\n", err.Error())
		os.Exit(1)
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop() // 程序结束时停止 ticker

	for {
		select {
		case t := <-ticker.C:
			interfacenewIPs, err := getInterfaceIPs() // 获取网卡和对应的 IP 地址
			if err != nil {
				fmt.Println("执行任务时间:", t)
				_, _ = fmt.Fprintf(mylog.SystemLogger, "CheckNetChange 获取网卡和对应的 IP 地址:%s\n", err.Error())
				os.Exit(1)
			}
			for ifaceName, newIfaceInfo := range interfacenewIPs {
				IfaceInfo, ok := interfaceIPs[ifaceName]
				if ok {
					if ipsEqualUnordered(IfaceInfo.Ips, newIfaceInfo.Ips) == false {
						_, _ = fmt.Fprintf(mylog.SystemLogger, "CheckNetChange  %s IP 地址变化:%v %v\n", ifaceName, IfaceInfo.Ips, newIfaceInfo.Ips)
						os.Exit(1)
					}
				} else {
					_, _ = fmt.Fprintf(mylog.SystemLogger, "CheckNetChange 新的接口:%s\n", ifaceName)
					os.Exit(1)
				}
			}
			// 在此处调用您需要每分钟执行的函数
		}
	}
}
