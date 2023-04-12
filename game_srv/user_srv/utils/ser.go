package utils

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"user_srv/global"
)

func RegistAndHealthCheak(server *grpc.Server, port int) (*api.Client, string) {
	//grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Fatalf("[RegistAndHealthCheak]服务注册连接失败:%s", err.Error())
	}

	reg := new(api.AgentServiceRegistration)
	reg.Name = global.ServerConfig.ConsulInfo.Name //服务name
	serverID := uuid.NewString()
	reg.ID = serverID //服务id
	reg.Port = port
	reg.Address = "192.168.159.1" //消费者访问服务地址
	//reg.Check = &api.AgentServiceCheck{
	//	GRPC: fmt.Sprintf("%s:%d", global.ServerConfig.Host, port),
	//}

	err = client.Agent().ServiceRegister(reg)
	if err != nil {
		zap.S().Fatalf("[RegistAndHealthCheak]服务注册失败:%s", err.Error())
	}
	return client, serverID
}
