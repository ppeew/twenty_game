package initialize

import (
	"fmt"
	"process_web/global"
	game_proto "process_web/proto/game"
	"time"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	go func() {
		for {
			select {
			case <-time.After(time.Second * 30):
				if global.GameSrvClient == nil {
					consulInfo := global.ServerConfig.ConsulInfo
					gameConn, err := grpc.Dial(
						fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GameSrvInfo.Name),
						grpc.WithInsecure(),
						grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
					)
					if err != nil {
						zap.S().Infof("[InitSrvConn] 连接 【游戏服务失败】")
					}
					gameSrvClient := game_proto.NewGameClient(gameConn)
					global.GameSrvClient = gameSrvClient
				}
			}
		}
	}()
}
