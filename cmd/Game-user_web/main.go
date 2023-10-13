package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"user_web/global"
	"user_web/initialize"
	"user_web/utils"

	"go.uber.org/zap"
)

func main() {
	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitDB()
	initialize.InitTasks()
	initialize.InitSrvConn()
	initialize.InitSentinel()

	routers := initialize.InitRouters()

	//自动找可用端口
	port, err := utils.GetFreePort()
	if err != nil {
		panic(err)
	}
	//port = 9000
	global.ServerConfig.Port = port
	if global.DEBUG {
		//是debug
		global.ServerConfig.Port = 9005
	}

	go func() {
		if err := routers.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
	}()
	//服务注册及健康检查
	consulClient, serverID := utils.RegistAndHealthCheck()
	//终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Kill, os.Interrupt)
	sig := <-quit
	zap.S().Infof("接收到退出信号 %+v\n", sig)

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
