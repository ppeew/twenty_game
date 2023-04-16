package global

import (
	"game_web/config"
	"game_web/proto"
	"github.com/redis/go-redis/v9"
)

var (
	RedisDB      *redis.Client
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient proto.GameClient
)
