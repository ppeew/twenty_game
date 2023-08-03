package main

import (
	"context"
	"fmt"
	"game_web/api"
	"game_web/global"
	"game_web/initialize"
	game_proto "game_web/proto/game"
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
	initialize.GetConsulServer()
	routers := initialize.InitRouters()
	utils.CheckGoRoutines()

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
	//资源释放 释放房间
	for roomID, _ := range api.CHAN {
		// 1.删除用户对应服务器连接

		// 2.删除redis房间信息
		global.GameSrvClient.DelRoomServer(context.Background(), &game_proto.RoomIDInfo{RoomID: roomID})
		// 3.删除该房间对应的服务器信息
		global.GameSrvClient.DeleteRoom(context.Background(), &game_proto.RoomIDInfo{RoomID: roomID})
	}
	zap.S().Info("释放资源完毕，退出")
}
