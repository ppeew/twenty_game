package initialize

import (
	"fmt"
	"game_web/global"
	game_proto "game_web/proto/game"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	gameConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GameSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【游戏服务失败】")
		return
	}
	gameSrvClient := game_proto.NewGameClient(gameConn)
	global.GameSrvClient = gameSrvClient
}
