package global

import (
	"go.mongodb.org/mongo-driver/mongo"
	"user_web/config"
	"user_web/proto/user"
)

var (
	DEBUG        bool
	MongoDB      *mongo.Database
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient user.UserClient
)
