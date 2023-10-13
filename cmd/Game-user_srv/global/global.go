package global

import (
	"go.mongodb.org/mongo-driver/mongo"
	"user_srv/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DEBUG        bool
	MysqlDB      *gorm.DB
	RedisDB      *redis.Client
	MongoDB      *mongo.Database
	NacosConfig  = &config.NacosConfig{}
	ServerConfig = &config.ServerConfig{}
	//GameSrvClient game.GameClient
)
