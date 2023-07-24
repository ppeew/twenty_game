package server

import (
	"context"
	"net/http"
	"process_web/global"
	"process_web/my_struct"
	"process_web/my_struct/response"
	"process_web/proto/game"
	"process_web/utils"
	"strconv"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	RoomQuit = iota
	GameStart
)

type RoomStruct struct {
	MsgChan  chan my_struct.Message //接受信息管道
	ExitChan chan int               //用于结束房间协程
	wg       sync.WaitGroup         //协调所有协程关闭

	RoomID        uint32
	MaxUserNumber uint32
	GameCount     uint32
	UserNumber    uint32
	RoomOwner     uint32
	RoomWait      bool
	Users         map[uint32]my_struct.UserRoomData
	RoomName      string
}

func NewRoomStruct(data *Data) RoomStruct {
	users := make(map[uint32]my_struct.UserRoomData)
	for _, userID := range data.users {
		//查询API用户信息
		var res utils.UserInfo
		gorequest.New().Get("http://139.159.234.134:8000/user/v1/search").Param("id", strconv.Itoa(int(userID))).
			Retry(5, time.Second, http.StatusInternalServerError).EndStruct(&res)
		users[userID] = my_struct.UserRoomData{
			ID:           userID,
			Ready:        false,
			IntoRoomTime: time.Now(),
			Nickname:     res.Nickname,
			Gender:       res.Gender,
			Username:     res.Username,
			Image:        res.Image,
		}
	}
	return RoomStruct{
		MsgChan:  make(chan my_struct.Message, 1024),
		ExitChan: make(chan int, 3),
		wg:       sync.WaitGroup{},

		RoomID:        data.roomID,
		MaxUserNumber: data.maxUserNumber,
		GameCount:     data.gameCount,
		UserNumber:    data.userNumber,
		RoomOwner:     data.roomOwner,
		RoomWait:      true,
		Users:         users,
		RoomName:      data.roomName,
	}
}

func (room *RoomStruct) RunRoom() (*Data, bool) {
	ctx, cancel := context.WithCancel(context.Background())
	dealFunc := NewDealFunc(room)
	for _, userData := range room.Users {
		go room.ReadRoomUserMsg(ctx, userData.ID)
	}
	//用于用户连接
	go room.ForUserConn(ctx)
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
			cancel()
			room.wg.Wait()
			if msg == RoomQuit {
				global.GameSrvClient.DeleteRoom(context.Background(), &game.RoomIDInfo{RoomID: room.RoomID})
				global.GameSrvClient.DelRoomServer(context.Background(), &game.RoomIDInfo{RoomID: room.RoomID})
				for id := range room.Users {
					global.GameSrvClient.DelConnData(context.Background(), &game.DelConnInfo{Id: id})
				}
				return nil, true
			} else if msg == GameStart {
				room.RoomWait = false
				var users []*game.RoomUser
				for _, data := range room.Users {
					users = append(users, &game.RoomUser{
						ID:    data.ID,
						Ready: data.Ready,
					})
				}
				global.GameSrvClient.SetGlobalRoom(context.Background(), &game.RoomInfo{
					RoomID:        room.RoomID,
					MaxUserNumber: room.MaxUserNumber,
					GameCount:     room.GameCount,
					UserNumber:    room.UserNumber,
					RoomOwner:     room.RoomOwner,
					RoomWait:      room.RoomWait,
					Users:         users,
					RoomName:      room.RoomName,
				})
				rsp := make([]uint32, 0)
				for userID := range room.Users {
					rsp = append(rsp, userID)
				}
				return NewData(room.RoomID, room.MaxUserNumber, room.GameCount, room.UserNumber, room.RoomOwner, room.RoomName, rsp), false
			}
		}
	}
}

// ReadRoomUserMsg 读取发送到房间的信息入管道
func (room *RoomStruct) ReadRoomUserMsg(ctx context.Context, userID uint32) {
	//当用户连接还没建立直接return，直到客户端调用连接
	if global.UsersConn[userID] == nil {
		//zap.S().Info("[ReadRoomUserMsg]]:用户连接没建立return")
		return
	}
	room.wg.Add(1)
	defer room.wg.Done()
	for true {
		select {
		case <-ctx.Done():
			return
		case message := <-global.UsersConn[userID].InChanRead():
			//zap.S().Infof("[ReadRoomUserMsg]:读到%d用户信息了", userID)
			message.UserID = userID //添加标识，能够识别用户
			room.MsgChan <- message
		case <-global.UsersConn[userID].IsDisConn():
			return
		}
	}
}

func (room *RoomStruct) CheckClientHealth(ctx context.Context) {
	room.wg.Add(1)
	defer room.wg.Done()
	for true {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(time.Second * 30):
			for _, info := range room.Users {
				if global.UsersConn[info.ID] != nil {
					err := global.UsersConn[info.ID].OutChanWrite(response.MessageResponse{MsgType: response.CheckHealthType})
					if err != nil {
						//检查用户连接，断开则自动离开房间
						room.MsgChan <- my_struct.Message{
							Type:         my_struct.QuitRoomMsg,
							UserID:       info.ID,
							QuitRoomData: my_struct.QuitRoomData{},
						}
					}
				}
			}
		}
	}
}

func (room *RoomStruct) ForUserConn(ctx context.Context) {
	room.wg.Add(1)
	defer room.wg.Done()
	for true {
		select {
		case userID := <-global.ConnectCHAN[room.RoomID]:
			//TODO 可能会出现并发问题 因此采用单线程处理
			go room.ReadRoomUserMsg(ctx, userID)
		case <-ctx.Done():
			return
		}
	}
}

func (room *RoomStruct) UpdateRedisRoom(ctx context.Context) {
	room.wg.Add(1)
	defer room.wg.Done()
	for true {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(time.Second * 3):
			var users []*game.RoomUser
			for _, data := range room.Users {
				users = append(users, &game.RoomUser{
					ID:    data.ID,
					Ready: data.Ready,
				})
			}
			global.GameSrvClient.SetGlobalRoom(context.Background(), &game.RoomInfo{
				RoomID:        room.RoomID,
				MaxUserNumber: room.MaxUserNumber,
				GameCount:     room.GameCount,
				UserNumber:    room.UserNumber,
				RoomOwner:     room.RoomOwner,
				RoomWait:      room.RoomWait,
				Users:         users,
				RoomName:      room.RoomName,
			})

		}
	}
}

func (room *RoomStruct) ForUserIntoRoom(ctx context.Context) {
	room.wg.Add(1)
	defer room.wg.Done()
	if global.IntoRoomCHAN[room.RoomID] == nil {
		global.IntoRoomCHAN[room.RoomID] = make(chan uint32)
	}
	for true {
		select {
		case <-ctx.Done():
			return
		case userID := <-global.IntoRoomCHAN[room.RoomID]:
			room.MsgChan <- my_struct.Message{Type: my_struct.UserIntoMsg, UserID: userID, UserIntoData: my_struct.UserIntoData{}}
		}
	}
}

func BroadcastToAllRoomUsers(roomInfo *RoomStruct, message response.MessageResponse) {
	for _, info := range roomInfo.Users {
		if global.UsersConn[info.ID] != nil {
			err := global.UsersConn[info.ID].OutChanWrite(message)
			if err != nil {
				//zap.S().Infof("ID为%d的用户掉线了", info.ID)
			}
		}
	}
}
