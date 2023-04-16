package main

import (
	"fmt"
	"game_web/global"
	"game_web/initialize"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitSrvConn()
	initialize.InitRedis()

	routers := initialize.InitRouters()

	if err := routers.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}

	//终止信号
	fmt.Println("okok!!!!!")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	//if err = register_client.DeRegister(serviceId); err != nil {
	//	zap.S().Info("注销失败:", err.Error())
	//}else{
	//	zap.S().Info("注销成功:")
	//}
}
