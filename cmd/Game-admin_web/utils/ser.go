package utils

import (
	"admin_web/global"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

func RegistAndHealthCheck() (*api.Client, string) {
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
	reg.Port = global.ServerConfig.Port
	reg.Address = global.ServerConfig.Host //消费者访问服务地址
	if !global.DEBUG {
		reg.Check = &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", global.ServerConfig.Host, global.ServerConfig.Port),
			Status:                         api.HealthPassing,
			DeregisterCriticalServiceAfter: "100s",
			Interval:                       "30s",
		}
	}
	err = client.Agent().ServiceRegister(reg)
	if err != nil {
		zap.S().Fatalf("[RegistAndHealthCheak]服务注册失败:%s", err.Error())
	}
	return client, serverID
}
