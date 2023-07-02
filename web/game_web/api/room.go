package api

import (
	"context"
	"game_web/global"
	"game_web/model"
	"game_web/model/response"
	"game_web/proto/game"
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

// 房间主函数
func startRoomThread(data RoomData) {
	ctx, cancel := context.WithCancel(context.Background())
	room := NewRoom(data)
	dealFunc := NewDealFunc(room)
	//读取房间内的管道 (正常来说，用户进入房间但是还没建立socket，此时连接为nil,该读取协程会关闭，当用户游戏结束，连接不为nil)
	for _, userData := range room.RoomData.Users {
		go room.ReadRoomUserMsg(ctx, userData.ID)
	}
	//用于用户进房
	go room.ForUserIntoRoom(ctx)
	//定时检查房间用户是否占用房间不退出（看socket是否断开了）
	go room.CheckClientHealth(ctx)
	//定时发送到redis，更新房间列表信息，为大厅外查询更新数据
	go room.UpdateRedisRoom(ctx)
	for {
		select {
		case msg := <-room.MsgChan:
			if dealFunc[msg.Type] != nil {
				dealFunc[msg.Type](msg)
			}
		case msg := <-room.ExitChan:
			// 停止信号，关闭主函数及相关子协程，优雅退出
			cancel()
			room.wg.Wait()
			zap.S().Info("[startRoomThread]]:其他协程已关闭")
			if msg == RoomQuit {
				global.GameSrvClient.DeleteRoom(context.Background(), &game.RoomIDInfo{RoomID: room.RoomData.RoomID})
				return
			} else if msg == GameStart {
				go RunGame(GameData{
					RoomID:     room.RoomData.RoomID,
					GameCount:  room.RoomData.GameCount,
					UserNumber: room.RoomData.MaxUserNumber,
					RoomOwner:  room.RoomData.RoomOwner,
					Users:      room.RoomData.Users,
					RoomName:   room.RoomData.RoomName,
				})
				return
			}
		}
	}
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
	if UsersConn[userID] == nil {
		return
	}
	roomInfo.wg.Add(1)
	defer roomInfo.wg.Done()
	for true {
		//fmt.Printf("[ReadRoomUserMsg] %+v,%+v,%+v\n", UsersConn, userID, UsersConn[userID])
		select {
		case <-ctx.Done():
			return
		case message := <-UsersConn[userID].InChanRead():
			message.UserID = userID //添加标识，能够识别用户
			roomInfo.MsgChan <- message
		case <-UsersConn[userID].IsDisConn():
			//zap.S().Infof("[ReadRoomUserMsg]:%d用户掉线了", userID)
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
			BroadcastToAllRoomUsers(roomInfo, response.MessageResponse{MsgType: response.CheckHealthType})
		}
	}
}

func (roomInfo *RoomStruct) ForUserIntoRoom(ctx context.Context) {
	roomInfo.wg.Add(1)
	defer roomInfo.wg.Done()
	for true {
		select {
		case userID := <-CHAN[roomInfo.RoomData.RoomID]:
			//TODO 可能会出现并发问题 因此采用单线程处理
			roomInfo.MsgChan <- model.Message{Type: model.UserIntoMsg, UserID: userID, UserIntoData: model.UserIntoData{UserID: userID}}
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
		case <-time.Tick(time.Second * 5):
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

func BroadcastToAllRoomUsers(roomInfo *RoomStruct, message response.MessageResponse) {
	for _, info := range roomInfo.RoomData.Users {
		if UsersConn[info.ID] != nil {
			err := UsersConn[info.ID].OutChanWrite(message)
			if err != nil {
				zap.S().Infof("ID为%d的用户掉线了", info.ID)
			}
		}
	}
}
