package main

import (
	"fmt"
	"os"
	"os/signal"
	"store_web/global"
	"store_web/initialize"
	"store_web/utils"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	routers := initialize.InitRouters()
	//utils.CreateTable()
	//自动获取可用端口号
	port, err := utils.GetFreePort()
	if err != nil {
		panic(err)
	}
	global.ServerConfig.Port = port

	go func() {
		if err := routers.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
	}()
	//服务注册及健康检查
	consulClient, serverID := utils.RegistAndHealthCheck()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Kill, os.Interrupt)
	<-quit
	// 资源释放
	go func() {
		<-quit
		zap.S().Info("两次ctrl+c强制退出")
		syscall.Exit(0)
	}()

	if err := consulClient.Agent().ServiceDeregister(serverID); err != nil {
		zap.S().Info("注销服务失败")
	}
	zap.S().Info("释放资源完毕，退出")
}
