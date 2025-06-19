package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"platform-backend/config"
	db "platform-backend/db"
	"platform-backend/migration"
	"platform-backend/routers"
	"platform-backend/service/message"
	"platform-backend/service/permission"
	"platform-backend/service/superadmin"
	"platform-backend/task"
	mylog "platform-backend/utils/log"

	_ "net/http/pprof"
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

	// 监听配置变更
	go config.DynamicReloadConfig()

	// 初始化log处理
	mylog.Init(config.MyConfig.AppConfig.Mod, config.MyConfig.AppConfig.LogExpire)

	//初始化db
	db.InitDB()
	if config.MyConfig.AutoMigration {
		migration.AutoMigration()
	}

	if err := permission.Migration(); err != nil {
		log.Panicf("初始化权限失败:%v", err)
		panic(err)
	}

	//初始化super admin账号
	if err := superadmin.InitSuperAdmin(); err != nil {
		log.Panicf("初始化super admin账号失败:%v", err)
	}

	//启动web服务
	go func() {
		if err := routers.Router.Init().Run(config.MyConfig.AppConfig.Address); err != nil {
			_, _ = fmt.Fprintf(mylog.SystemLogger, "启动web服务失败:%s\n", err.Error())
			errChan <- err
		}
	}()

	//启动定时任务
	go func() {
		if err := task.RunTask(); err != nil {
			_, _ = fmt.Fprintf(mylog.SystemLogger, "启动任务失败:%s\n", err.Error())
			errChan <- err
		}
	}()

	//server, err := protocol.StartMdns()
	//if err != nil {
	//	errChan <- err
	//} else {
	//	defer server.Shutdown()
	//}

	//启动mqtt并开始订阅
	err := message.Start()
	if err != nil {
		log.Panicf("启动mqtt失败:%v", err)
	}

	if err := <-errChan; err != nil {
		log.Fatal(err)
	}
}
