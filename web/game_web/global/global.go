package global

import (
	"game_web/config"
	game_proto "game_web/proto/game"
)

var (
	DEBUG         bool
	NacosConfig   *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	GameSrvClient game_proto.GameClient
	//UserSrvClient user_proto.UserClient
)
