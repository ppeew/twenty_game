package global

import (
	"game_srv/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	MysqlDB      *gorm.DB
	RedisDB      *redis.Client
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
