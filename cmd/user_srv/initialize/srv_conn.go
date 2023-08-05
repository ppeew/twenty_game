package initialize

import (
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
)

//func InitSrvConn() {
//	consulInfo := global.ServerConfig.ConsulInfo
//	gameConn, err := grpc.Dial(
//		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GameSrvInfo.Name),
//		grpc.WithInsecure(),
//		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
//	)
//	if err != nil {
//		zap.S().Fatal("[InitSrvConn] 连接 【游戏物品服务失败】")
//	}
//	global.GameSrvClient = game.NewGameClient(gameConn)
//}
