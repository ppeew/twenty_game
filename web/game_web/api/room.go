package api

import (
	"context"
	"fmt"
	"game_web/global"
	"game_web/model"
	"game_web/model/response"
	game_proto "game_web/proto/game"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Room struct {
	RoomID   uint32
	MsgChan  chan model.Message //接受信息管道
	ExitChan chan int           //用于结束房间协程
	wg       sync.WaitGroup     //协调所有协程关闭
}

func NewRoom(roomID uint32) *Room {
	room := &Room{
		RoomID:   roomID,
		MsgChan:  make(chan model.Message, 1024),
		ExitChan: make(chan int, 3),
		wg:       sync.WaitGroup{},
	}
	CHAN[roomID] = make(chan uint32, 10)
	return room
}

// 房间主函数
func startRoomThread(roomID uint32) {
	//初始化房间信息
	room := NewRoom(roomID)
	dealFunc := NewDealFunc(room)
	ctx, cancel := context.WithCancel(context.Background())
	searchRoom, err := global.GameSrvClient.SearchRoom(context.Background(), &game_proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		//zap.S().Infof("[startRoomThread]:%s", err)
		cancel()
		return
	}
	//BroadcastToAllRoomUsers(searchRoom, GrpcModelToResponse(searchRoom))
	for _, info := range searchRoom.Users {
		go room.ReadRoomUserMsg(ctx, info.ID)
		room.wg.Add(1)
	}
	go func(ctx context.Context) {
		room.wg.Add(1)
		for true {
			select {
			case userID := <-CHAN[roomID]:
				go room.ReadRoomUserMsg(ctx, userID)
				room.wg.Add(1)
			case <-ctx.Done():
				room.wg.Done()
				return
			}
		}
	}(ctx)
	for {
		select {
		case msg := <-room.MsgChan:
			//zap.S().Infof("收到%+v", msg)
			if msg.Type >= model.QuitRoomMsg && msg.Type <= model.RoomBeginGameMsg {
				//只有这类消息才处理
				dealFunc[msg.Type](msg)
			}
		case msg := <-room.ExitChan:
			// 停止信号，关闭主函数及读取用户通道函数，优雅退出
			cancel()
			//zap.S().Info("[startRoomThread]]:wg正在等待其他协程结束")
			room.wg.Wait()
			//zap.S().Info("[startRoomThread]]:其他协程已关闭")
			if msg == model.RoomQuit {
				return
			} else if msg == model.GameStart {
				go RunGame(roomID)
				return
			}
		}
	}
}

// ReadRoomUserMsg 读取发送到房间的信息入管道
func (roomInfo *Room) ReadRoomUserMsg(ctx context.Context, userID uint32) {
	for true {
		select {
		case <-ctx.Done():
			roomInfo.wg.Done()
			return
		case message := <-UsersState[userID].InChanRead():
			message.UserID = userID //添加标识，能够识别用户
			roomInfo.MsgChan <- message
		case <-UsersState[userID].IsDisConn():
			//如果与用户的websocket断开，退出读取协程,并且将该玩家从房间剔除
			zap.S().Infof("[ReadRoomUserMsg]:%d用户掉线了", userID)
			roomInfo.wg.Done()
			roomInfo.MsgChan <- model.Message{
				Type:   model.QuitRoomMsg,
				UserID: userID,
			}
			return
		}
	}
}

// RoomInfo 房间信息
func (roomInfo *Room) RoomInfo(message model.Message) {
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &game_proto.RoomIDInfo{RoomID: roomInfo.RoomID})
	if err != nil {
		zap.S().Error("[RoomInfo]:%s", err)
		return
	}
	resp := GrpcModelToResponse(room)
	SendMsgToUser(UsersState[message.UserID], response.MessageResponse{
		MsgType:  response.RoomInfoResponseType,
		RoomInfo: resp,
	})
}

// QuitRoom 退出房间（房主退出会导致全部房间删除）
func (roomInfo *Room) QuitRoom(message model.Message) {
	info, err := global.GameSrvClient.QuitRoom(context.Background(), &game_proto.QuitRoomInfo{
		RoomID: roomInfo.RoomID,
		UserID: message.UserID,
	})
	if err != nil {
		zap.S().Warnf("[QuitRoom]:%s", err)
		return
	}
	if info.IsOwnerQuit {
		// 房主退出会销毁房间
		SendMsgToUser(UsersState[message.UserID], response.MessageResponse{
			MsgType: response.MsgResponseType,
			MsgInfo: response.MsgResponse{
				MsgData: "房主退出房间成功",
			}})
		BroadcastToAllRoomUsers(info.RoomInfo, response.MessageResponse{
			MsgType: response.MsgResponseType,
			MsgInfo: response.MsgResponse{
				MsgData: "房主退出房间，房间将关闭",
			},
		})
		time.Sleep(1 * time.Second)
		UsersState[message.UserID].CloseConn()
		for _, info := range info.RoomInfo.Users {
			UsersState[info.ID].CloseConn()
		}
		roomInfo.ExitChan <- model.RoomQuit
	} else {
		//游戏玩家的退出
		UsersState[message.UserID].CloseConn()
		//房间变化，广播
		resp := GrpcModelToResponse(info.RoomInfo)
		BroadcastToAllRoomUsers(info.RoomInfo, response.MessageResponse{
			MsgType:  response.RoomInfoResponseType,
			RoomInfo: resp,
		})
	}
}

// UpdateRoom 更新房间的房主或者游戏配置(仅房主)
func (roomInfo *Room) UpdateRoom(message model.Message) {
	room, err := global.GameSrvClient.UpdateRoom(context.Background(), &game_proto.UpdateRoomInfo{
		RoomID:        roomInfo.RoomID,
		MaxUserNumber: message.UpdateData.MaxUserNumber,
		GameCount:     message.UpdateData.GameCount,
		Owner:         message.UpdateData.Owner,
		Kicker:        message.UpdateData.Kicker,
	})
	if err != nil {
		zap.S().Error("[UpdateRoom]:%s", err)
		return
	}
	if message.UpdateData.Kicker != 0 {
		//发送给被t的玩家
		SendMsgToUser(UsersState[message.UpdateData.Kicker], response.MessageResponse{
			MsgType: response.KickerResponseType, KickerInfo: response.KickerResponse{},
		})
		//t人了还需要关闭房间里面的连接(等待一段时间再关闭连接，为了被t的玩家已经收到被t信息了)
		time.Sleep(1 * time.Second)
		UsersState[message.UpdateData.Kicker].CloseConn()
	}
	SendMsgToUser(UsersState[message.UserID], response.MessageResponse{
		MsgType: response.MsgResponseType,
		MsgInfo: response.MsgResponse{
			MsgData: "更新房间成功",
		},
	})
	//更新房间，发送广播
	resp := GrpcModelToResponse(room)
	BroadcastToAllRoomUsers(room, response.MessageResponse{
		MsgType:  response.RoomInfoResponseType,
		RoomInfo: resp,
	})
}

// UpdateUserReadyState 玩家准备状态
func (roomInfo *Room) UpdateUserReadyState(message model.Message) {
	room, err := global.GameSrvClient.UpdateUserReadyState(context.Background(), &game_proto.ReadyStateInfo{
		RoomID:  roomInfo.RoomID,
		UserID:  message.UserID,
		IsReady: message.ReadyStateData.IsReady,
	})
	if err != nil {
		zap.S().Infof("[UpdateUserReadyState]:%s", err.Error())
		return
	}
	SendMsgToUser(UsersState[message.UserID], response.MessageResponse{
		MsgType: response.MsgResponseType,
		MsgInfo: response.MsgResponse{
			MsgData: fmt.Sprintf("玩家%d准备状态更新", message.UserID),
		},
	})
	resp := GrpcModelToResponse(room)
	BroadcastToAllRoomUsers(room, response.MessageResponse{
		MsgType:  response.RoomInfoResponseType,
		RoomInfo: resp,
	})
}

// BeginGame 开始游戏
func (roomInfo *Room) BeginGame(message model.Message) {
	room, err := global.GameSrvClient.BeginGame(context.Background(), &game_proto.BeginGameInfo{
		RoomID: roomInfo.RoomID,
		UserID: message.UserID,
	})
	if err != nil {
		zap.S().Infof("[BeginGame]:%s", err)
		return
	}
	if room.ErrorMsg != "" {
		SendMsgToUser(UsersState[message.UserID], response.MessageResponse{
			MsgType: response.MsgResponseType,
			MsgInfo: response.MsgResponse{
				MsgData: room.ErrorMsg,
			},
		})
		return
	}
	roomInfo.ExitChan <- model.GameStart
}

func (roomInfo *Room) ChatProcess(message model.Message) {
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &game_proto.RoomIDInfo{RoomID: roomInfo.RoomID})
	if err != nil {
		zap.S().Errorf("[ChatProcess]:%s", err)
	}
	BroadcastToAllRoomUsers(room, response.MessageResponse{
		MsgType: response.MsgResponseType,
		MsgInfo: response.MsgResponse{
			MsgData: fmt.Sprintf("用户%d说：%s", message.UserID, string(message.ChatMsgData.Data)),
		},
	})
}

func (roomInfo *Room) CheckHealth(message model.Message) {
	SendMsgToUser(UsersState[message.UserID], response.MessageResponse{
		MsgType:    response.CheckHealthResponseType,
		HealthInfo: response.HealthResponse{},
	})
}

func GrpcModelToResponse(roomInfo *game_proto.RoomInfo) response.RoomResponse {
	resp := response.RoomResponse{
		RoomID:        roomInfo.RoomID,
		MaxUserNumber: roomInfo.MaxUserNumber,
		GameCount:     roomInfo.GameCount,
		UserNumber:    roomInfo.UserNumber,
		RoomOwner:     roomInfo.RoomOwner,
		RoomWait:      roomInfo.RoomWait,
		RoomName:      roomInfo.RoomName,
	}
	for _, roomUser := range roomInfo.Users {
		resp.Users = append(resp.Users, response.UserResponse{
			ID:    roomUser.ID,
			Ready: roomUser.Ready,
		})
	}
	return resp
}

func BroadcastToAllRoomUsers(roomInfo *game_proto.RoomInfo, message response.MessageResponse) {
	for _, info := range roomInfo.Users {
		err := UsersState[info.ID].OutChanWrite(message)
		if err != nil {
			zap.S().Infof("ID为%d的用户掉线了", info.ID)
		}
	}
}

type dealFunc func(message model.Message)

func NewDealFunc(room *Room) map[uint32]dealFunc {
	var dealFun = make(map[uint32]dealFunc)
	dealFun[model.CheckHealthMsg] = room.CheckHealth
	dealFun[model.QuitRoomMsg] = room.QuitRoom
	dealFun[model.GetRoomMsg] = room.RoomInfo
	dealFun[model.RoomBeginGameMsg] = room.BeginGame
	dealFun[model.UserReadyStateMsg] = room.UpdateUserReadyState
	dealFun[model.UpdateRoomMsg] = room.UpdateRoom
	dealFun[model.ChatMsg] = room.ChatProcess
	return dealFun
}
