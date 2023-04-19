package model

type RoomInfo struct {
	RoomID        uint32
	MaxUserNumber uint32
	GameCount     uint32
	UserNumber    uint32
	RoomOwner     uint32
	RoomWait      bool
	//存储用户连接相关
	Publisher *Publisher
}
