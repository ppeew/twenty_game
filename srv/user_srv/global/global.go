package global

import (
	"user_srv/config"
	"user_srv/proto/game"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DEBUG         bool
	MysqlDB       *gorm.DB
	RedisDB       *redis.Client
	NacosConfig   = &config.NacosConfig{}
	ServerConfig  = &config.ServerConfig{}
	GameSrvClient game.GameClient
)
