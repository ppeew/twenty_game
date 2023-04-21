package global

import (
	"game_web/config"
	"game_web/proto"
)

var (
	NacosConfig   *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	GameSrvClient proto.GameClient
)
