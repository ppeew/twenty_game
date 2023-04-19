package global

import (
	"game_web/config"
	"game_web/model"
	"game_web/proto"
)

var (
	NacosConfig   *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	GameSrvClient proto.GameClient

	//游戏相关
	RoomData map[uint32]*model.RoomInfo = make(map[uint32]*model.RoomInfo) //房间号->房间数据的映射(每个房间线程访问各自数据)
)
