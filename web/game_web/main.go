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
	routers := initialize.InitRouters()
	go func() {
		if err := routers.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
		zap.S().Info("启动成功")
	}()
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
	//for roomID, _ := range api.RoomData {
	//	_, err := global.GameSrvClient.DeleteRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	//	if err != nil {
	//		zap.S().Infof("关闭%d房间失败:%s", roomID, err.Error())
	//	}
	//	zap.S().Infof("已经将%d房间清空！", roomID)
	//}
	zap.S().Info("释放资源完毕，退出")
	//if err = register_client.DeRegister(serviceId); err != nil {
	//	zap.S().Info("注销失败:", err.Error())
	//}else{
	//	zap.S().Info("注销成功:")
	//}
}
