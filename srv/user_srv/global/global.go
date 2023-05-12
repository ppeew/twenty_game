package global

import (
	"user_srv/config"
	"user_srv/proto/game"

	"gorm.io/gorm"
)

var (
	DEBUG         bool
	DB            *gorm.DB
	NacosConfig   = &config.NacosConfig{}
	ServerConfig  = &config.ServerConfig{}
	GameSrvClient game.GameClient
)
