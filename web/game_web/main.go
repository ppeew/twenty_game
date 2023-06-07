package main

import (
	"fmt"
	"game_web/global"
	"game_web/initialize"
	"game_web/utils"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitSrvConn()
	initialize.InitSentinel()
	routers := initialize.InitRouters()
	//utils.CheckGoRoutines()

	//服务注册及健康检查
	consulClient, serverID := utils.RegistAndHealthCheck()
	go func() {
		if err := routers.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
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
	//资源释放
	utils.ReleaseResource()
	zap.S().Info("释放资源完毕，退出")
}
