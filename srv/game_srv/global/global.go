package global

import (
	"game_srv/config"
	"game_srv/proto/user"

	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DEBUG         bool
	MysqlDB       *gorm.DB
	RedisDB       *redis.Client
	RedSync       *redsync.Redsync
	NacosConfig   *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	UserSrvClient user.UserClient
)
