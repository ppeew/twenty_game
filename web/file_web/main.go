package main

import (
	"file_web/global"
	"file_web/initialize"
	"file_web/utils"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitSrvConn()

	routers := initialize.InitRouters()
	//自动找可用端口
	port, err := utils.GetFreePort()
	if err != nil {
		panic(err)
	}
	global.ServerConfig.Port = port

	go func() {
		if err := routers.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
		zap.S().Info("启动成功")
	}()
	consulClient, serverID := utils.RegistAndHealthCheck()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit
	// 资源释放
	go func() {
		// 强制退出
		<-quit
		zap.S().Info("两次ctrl+c强制退出")
		syscall.Exit(0)
	}()
	if err := consulClient.Agent().ServiceDeregister(serverID); err != nil {
		zap.S().Info("注销服务失败")
	}
	//if err = register_client.DeRegister(serviceId); err != nil {
	//	zap.S().Info("注销失败:", err.Error())
	//}else{
	//	zap.S().Info("注销成功:")
	//}
}
