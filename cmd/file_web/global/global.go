package global

import (
	"file_web/config"
	"file_web/proto/user"
)

var (
	DEBUG        bool
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient user.UserClient
)
