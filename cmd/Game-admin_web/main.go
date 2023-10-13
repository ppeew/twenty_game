package main

import (
	"admin_web/global"
	"admin_web/initialize"
	"admin_web/utils"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	// 先获取配置才能初始化logger,db,task
	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitDB()
	initialize.InitTasks()
	go initialize.GetConsulServer()
	routers := initialize.InitRouters()
	//utils.CreateTable()
	//自动获取可用端口号
	port, err := utils.GetFreePort()
	if err != nil {
		panic(err)
	}
	global.ServerConfig.Port = port
	if global.DEBUG {
		//是debug
		global.ServerConfig.Port = 9100
	}
	go func() {
		if err := routers.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
	}()
	//服务注册及健康检查
	consulClient, serverID := utils.RegistAndHealthCheck()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Kill, os.Interrupt)
	sig := <-quit
	zap.S().Infof("接收到退出信号 %+v\n", sig)

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
