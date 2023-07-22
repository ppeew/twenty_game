package global

import (
	"user_web/config"
	"user_web/proto/user"
)

var (
	DEBUG        bool
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient user.UserClient
)
