package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"process_web/global"
	"process_web/initialize"
	game_proto "process_web/proto/game"
	"process_web/utils"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitDB()
	initialize.InitTasks()
	initialize.InitSrvConn()
	routers := initialize.InitRouters()
	utils.CheckGoRoutines()

	//自动获取可用端口号
	port, err := utils.GetFreePort()
	if err != nil {
		panic(err)
	}
	global.ServerConfig.Port = port
	if global.DEBUG {
		//是debug
		global.ServerConfig.Port = 9002
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
	//资源释放 释放房间
	// 1.删除用户对应服务器连接
	//for id := range global.UsersConn {
	//	global.GameSrvClient.DelConnData(context.Background(), &game_proto.DelConnInfo{Id: id})
	//}
	global.UsersConn.Range(func(key, value any) bool {
		global.GameSrvClient.DelConnData(context.Background(), &game_proto.DelConnInfo{Id: key.(uint32)})
		return true
	})
	for roomID, _ := range global.ConnectCHAN {
		// 2.删除redis房间信息
		global.GameSrvClient.DelRoomServer(context.Background(), &game_proto.RoomIDInfo{RoomID: roomID})
		// 3.删除该房间对应的服务器信息
		global.GameSrvClient.DeleteRoom(context.Background(), &game_proto.RoomIDInfo{RoomID: roomID})
	}
	zap.S().Info("释放资源完毕，退出")
}
