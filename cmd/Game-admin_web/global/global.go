package global

import (
	"admin_web/config"
	"github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DEBUG                 bool
	MysqlDB               *gorm.DB
	MongoDB               *mongo.Database
	RedisDB               *redis.Client
	RedSync               *redsync.Redsync
	NacosConfig           *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig          *config.ServerConfig = &config.ServerConfig{}
	ConsulHallWebServices []*api.CatalogService
)
