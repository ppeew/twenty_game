package global

import (
	"github.com/hashicorp/consul/api"
	"hall_web/config"

	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DEBUG                 bool
	MysqlDB               *gorm.DB
	RedisDB               *redis.Client
	RedSync               *redsync.Redsync
	NacosConfig           *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig          *config.ServerConfig = &config.ServerConfig{}
	ConsulHallWebServices []*api.CatalogService
)
