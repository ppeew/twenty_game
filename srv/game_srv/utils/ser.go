package utils

import (
	"fmt"
	"game_srv/global"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RegistAndHealthCheck(server *grpc.Server, port int) (*api.Client, string) {
	if !global.DEBUG {
		grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	}

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Fatalf("[RegistAndHealthCheak]服务注册连接失败:%s", err.Error())
	}

	reg := new(api.AgentServiceRegistration)
	reg.Name = global.ServerConfig.ConsulInfo.Name //服务name
	serverID := uuid.NewString()
	//serverID := uuid.New().String()
	reg.ID = serverID //服务id
	reg.Port = port
	reg.Address = global.ServerConfig.ConsulInfo.ServerHost //消费者访问服务地址
	if !global.DEBUG {
		reg.Check = &api.AgentServiceCheck{
			GRPC: fmt.Sprintf("%s:%d", global.ServerConfig.Host, port),
		}
	}

	err = client.Agent().ServiceRegister(reg)
	if err != nil {
		zap.S().Fatalf("[RegistAndHealthCheak]服务注册失败:%s", err.Error())
	}
	return client, serverID
}
