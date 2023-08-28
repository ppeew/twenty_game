package global

import (
	"process_web/config"
	game_proto "process_web/proto/game"
	"sync"
)

// ConnectCHAN 房间号对应创建读取协程的管道
var ConnectCHAN = make(map[uint32]chan uint32)

// IntoRoomCHAN 用户进房发送chan 房间服务器读取并处理 key:房间号 value:用户id
var IntoRoomCHAN = make(map[uint32]chan uint32)

// IntoRoomRspCHAN IntoRoomChan 对用户进房做出回复  key:房间号 value:加入是否成功
var IntoRoomRspCHAN = make(map[uint32]chan bool)

// UsersConn 用户ID -> 用户连接
// var UsersConn = make(map[uint32]*WSConn)
var UsersConn = sync.Map{}

var (
	DEBUG         bool
	NacosConfig   *config.NacosConfig  = &config.NacosConfig{}
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	GameSrvClient game_proto.GameClient
)
