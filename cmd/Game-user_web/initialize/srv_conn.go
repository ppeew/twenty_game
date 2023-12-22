package initialize

import (
	"fmt"
	"time"
	"user_web/global"
	"user_web/proto/user"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	go func() {
		for {
			select {
			case <-time.After(time.Second * 30):
				if global.UserSrvClient == nil {
					consulInfo := global.ServerConfig.ConsulInfo
					userConn, err := grpc.Dial(
						fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
						grpc.WithInsecure(),
						grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
					)
					if err != nil {
						zap.S().Infof("[InitSrvConn] 连接 【用户服务失败】")
					}

					userSrvClient := user.NewUserClient(userConn)
					global.UserSrvClient = userSrvClient
				}
			}
		}
	}()
}
