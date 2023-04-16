package global

import (
	"gorm.io/gorm"
	"user_srv/config"
)

var (
	DB           *gorm.DB
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
