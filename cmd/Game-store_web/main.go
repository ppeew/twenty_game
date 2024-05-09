package main

import (
	"fmt"
	"os"
	"os/signal"
	"store_web/global"
	"store_web/initialize"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	initialize.InitConfig2()
	initialize.InitLogger()
	initialize.InitDB()
	initialize.InitTasks()
	routers := initialize.InitRouters()
	//utils.CreateTable()
	if global.DEBUG {
		//是debug
		global.ServerConfig.Port = 9001
	}
	go func() {
		if err := routers.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
	}()
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

	zap.S().Info("释放资源完毕，退出")
}
