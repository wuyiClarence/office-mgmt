package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"platform-mdns/config"
	"platform-mdns/mdns"
	mylog "platform-mdns/utils/log"
)

var (
	errChan = make(chan error)
)

func main() {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	if err := config.InitConfig(); err != nil {
		_, _ = fmt.Fprintf(mylog.SystemLogger, "初始化配置失败:%s\n", err.Error())
		errChan <- err
	}

	go mdns.CheckNetChange()

	// 监听配置变更
	go config.DynamicReloadConfig()

	mylog.Init(config.MyConfig.AppConfig.Mod, config.MyConfig.AppConfig.LogExpire)

	servers, err := mdns.StartMdns()
	if err != nil {
		errChan <- err
	}

	if err := <-errChan; err != nil {
		for _, server := range servers {
			_ = server.Shutdown()
		}

		log.Fatal(err)
	}
}
