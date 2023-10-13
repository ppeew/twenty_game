package global

import (
	"file_web/config"
	"file_web/proto/user"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DEBUG        bool
	MongoDB      *mongo.Database
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient user.UserClient
)
