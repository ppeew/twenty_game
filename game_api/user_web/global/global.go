package global

import (
	"user_web/config"
	"user_web/proto"
)

var (
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient proto.UserClient
)
