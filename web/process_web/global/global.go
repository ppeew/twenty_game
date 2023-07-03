package global

import (
	"process_web/config"
	game_proto "process_web/proto/game"
)

var (
	DEBUG         bool
	NacosConfig   *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	GameSrvClient game_proto.GameClient
)
