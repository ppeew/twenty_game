package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"user_srv/global"
	"user_srv/handler"
	"user_srv/initialize"
	"user_srv/proto/user"
	"user_srv/utils"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitDB()
	initialize.InitTasks()
	//initialize.InitSrvConn()
	server := grpc.NewServer()
	user.RegisterUserServer(server, &handler.UserServer{})

	//自动获取可用端口号
	//port, err := utils.GetFreePort()
	//if err != nil {
	//	zap.S().Fatalf("无法找到适用的端口号:%s", err)
	//}

	port := global.ServerConfig.Port
	zap.S().Infof("开启端口是:%d", port)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		zap.S().Errorf("启动listen失败:%s", err.Error())
	}

	//服务注册及健康检查
	consulClient, serverID := utils.RegistAndHealthCheck(server, port)

	//启动微服务
	go func() {
		err2 := server.Serve(listen)
		if err2 != nil {
			zap.S().Fatalf("启动grpc服务失败:%s", err2.Error())
		}
	}()
	zap.S().Info("启动服务成功")

	//优雅退出
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
	zap.S().Info("注销服务成功")
}
