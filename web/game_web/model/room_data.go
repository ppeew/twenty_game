package model

type RoomCoon struct {
	//存储用户连接相关
	MsgChan   chan Message       //接受信息管道
	UsersConn map[uint32]*WSConn //用户id到订阅者的映射
}
