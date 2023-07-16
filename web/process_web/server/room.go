package server

import (
	"context"
	"process_web/global"
	"process_web/model"
	"process_web/model/response"
	"process_web/proto/game"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	RoomQuit = iota
	GameStart
)

type RoomStruct struct {
	MsgChan  chan model.Message //接受信息管道
	ExitChan chan int           //用于结束房间协程
	wg       sync.WaitGroup     //协调所有协程关闭

	RoomData RoomData
}

type RoomData struct {
	RoomID        uint32
	MaxUserNumber uint32
	GameCount     uint32
	UserNumber    uint32
	RoomOwner     uint32
	RoomWait      bool
	Users         map[uint32]response.UserData
	RoomName      string
}

func NewRoom(data RoomData) *RoomStruct {
	room := &RoomStruct{
		MsgChan:  make(chan model.Message, 1024),
		ExitChan: make(chan int, 3),
		wg:       sync.WaitGroup{},
		RoomData: data,
	}
	return room
}

// ReadRoomUserMsg 读取发送到房间的信息入管道
func (roomInfo *RoomStruct) ReadRoomUserMsg(ctx context.Context, userID uint32) {
	//当用户连接还没建立直接return，直到客户端调用连接
	if global.UsersConn[userID] == nil {
		zap.S().Info("[ReadRoomUserMsg]]:用户连接没建立return")
		return
	}
	roomInfo.wg.Add(1)
	defer roomInfo.wg.Done()
	for true {
		//fmt.Printf("[ReadRoomUserMsg] %+v,%+v,%+v\n", global.UsersConn, userID, global.UsersConn[userID])
		select {
		case <-ctx.Done():
			return
		case message := <-global.UsersConn[userID].InChanRead():
			zap.S().Infof("[ReadRoomUserMsg]:读到%d用户信息了", userID)
			message.UserID = userID //添加标识，能够识别用户
			roomInfo.MsgChan <- message
		case <-global.UsersConn[userID].IsDisConn():
			zap.S().Infof("[ReadRoomUserMsg]:%d用户掉线了", userID)
			return
		}
	}
}

func (roomInfo *RoomStruct) CheckClientHealth(ctx context.Context) {
	roomInfo.wg.Add(1)
	defer roomInfo.wg.Done()
	for true {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(time.Second * 30):
			for _, info := range roomInfo.RoomData.Users {
				if global.UsersConn[info.ID] != nil {
					err := global.UsersConn[info.ID].OutChanWrite(response.MessageResponse{MsgType: response.CheckHealthType})
					if err != nil {
						//检查用户连接，断开则自动离开房间
						roomInfo.MsgChan <- model.Message{
							Type:         model.QuitRoomMsg,
							UserID:       info.ID,
							QuitRoomData: model.QuitRoomData{},
						}
					}
				}
			}
			//BroadcastToAllRoomUsers(roomInfo, response.MessageResponse{MsgType: response.CheckHealthType})
		}
	}
}

func (roomInfo *RoomStruct) ForUserConn(ctx context.Context) {
	roomInfo.wg.Add(1)
	defer roomInfo.wg.Done()
	for true {
		select {
		case userID := <-global.ConnectCHAN[roomInfo.RoomData.RoomID]:
			//TODO 可能会出现并发问题 因此采用单线程处理
			go roomInfo.ReadRoomUserMsg(ctx, userID)
		case <-ctx.Done():
			return
		}
	}
}

func (roomInfo *RoomStruct) UpdateRedisRoom(ctx context.Context) {
	roomInfo.wg.Add(1)
	defer roomInfo.wg.Done()
	for true {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(time.Second * 3):
			var users []*game.RoomUser
			for _, data := range roomInfo.RoomData.Users {
				users = append(users, &game.RoomUser{
					ID:    data.ID,
					Ready: data.Ready,
				})
			}
			global.GameSrvClient.SetGlobalRoom(context.Background(), &game.RoomInfo{
				RoomID:        roomInfo.RoomData.RoomID,
				MaxUserNumber: roomInfo.RoomData.MaxUserNumber,
				GameCount:     roomInfo.RoomData.GameCount,
				UserNumber:    roomInfo.RoomData.UserNumber,
				RoomOwner:     roomInfo.RoomData.RoomOwner,
				RoomWait:      roomInfo.RoomData.RoomWait,
				Users:         users,
				RoomName:      roomInfo.RoomData.RoomName,
			})

		}
	}
}

func (roomInfo *RoomStruct) ForUserIntoRoom(ctx context.Context) {
	roomInfo.wg.Add(1)
	defer roomInfo.wg.Done()
	if global.IntoRoomCHAN[roomInfo.RoomData.RoomID] == nil {
		global.IntoRoomCHAN[roomInfo.RoomData.RoomID] = make(chan uint32)
	}
	for true {
		select {
		case <-ctx.Done():
			return
		case userID := <-global.IntoRoomCHAN[roomInfo.RoomData.RoomID]:
			//读到用户进房消息
			//zap.S().Infof("[ForUserIntoRoom]:我看到你进房了，正在处理！")
			roomInfo.MsgChan <- model.Message{Type: model.UserIntoMsg, UserID: userID, UserIntoData: model.UserIntoData{}}
		}

	}
}

func BroadcastToAllRoomUsers(roomInfo *RoomStruct, message response.MessageResponse) {
	for _, info := range roomInfo.RoomData.Users {
		if global.UsersConn[info.ID] != nil {
			err := global.UsersConn[info.ID].OutChanWrite(message)
			if err != nil {
				//zap.S().Infof("ID为%d的用户掉线了", info.ID)
			}
		}
	}
}
