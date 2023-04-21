package model

type RoomCoon struct {
	//存储用户连接相关
	MsgChan   chan Message       //接受信息管道
	ExitChan  chan struct{}      //用于结束房间协程
	ReadExit  chan struct{}      //告知读用户消息线程是否已经完成退出
	UsersConn map[uint32]*WSConn //用户id到订阅者的映射
}
