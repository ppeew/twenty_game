package global

import (
	"user_srv/config"

	"gorm.io/gorm"
)

var (
	DEBUG        bool
	DB           *gorm.DB
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
