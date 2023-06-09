package main

import (
	"fmt"
	"game_srv/global"
	"game_srv/handler"
	"game_srv/initialize"
	"game_srv/proto/game"
	"game_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitSrvConn()

	server := grpc.NewServer()
	game.RegisterGameServer(server, &handler.GameServer{})

	port, err := utils.GetFreePort()
	global.ServerConfig.Port = port
	if err != nil {
		zap.S().Fatalf("无法找到适用的端口号:%s", err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		zap.S().Errorf("启动listen失败:%s", err.Error())
	}

	consulClient, serverID := utils.RegistAndHealthCheck(server, port)

	go func() {
		err2 := server.Serve(listen)
		if err2 != nil {
			zap.S().Fatalf("启动grpc服务失败:%s", err2.Error())
		}
	}()
	zap.S().Info("启动服务成功")
	//优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit
	if err := consulClient.Agent().ServiceDeregister(serverID); err != nil {
		zap.S().Info("注销服务失败")
	}
}
