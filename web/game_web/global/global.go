package global

import (
	"game_web/config"
	game_proto "game_web/proto/game"
)

var (
	DEBUG         bool
	NacosConfig   = &config.NacosConfig{}
	ServerConfig  = &config.ServerConfig{}
	GameSrvClient game_proto.GameClient
)
