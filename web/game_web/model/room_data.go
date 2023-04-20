package model

type RoomCoon struct {
	//存储用户连接相关
	MsgChan     chan Message       //接受信息管道
	PauseChan   chan struct{}      //用于暂停房间协程，等待游戏协程唤醒
	RecoverChan chan struct{}      //用于游戏结束恢复房间协程
	UsersConn   map[uint32]*WSConn //用户id到订阅者的映射
}
