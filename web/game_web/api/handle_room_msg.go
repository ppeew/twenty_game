package api

import (
	"context"
	"fmt"
	"game_web/global"
	"game_web/model"
	"game_web/model/response"
	game_proto "game_web/proto/game"
	"math/rand"
)

type dealFunc func(message model.Message)

func NewDealFunc(room *RoomStruct) map[uint32]dealFunc {
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

func (roomInfo *RoomStruct) MakeRoomResponse() *response.RoomResponse {
	roomResponse := &response.RoomResponse{
		RoomID:        roomInfo.RoomData.RoomID,
		MaxUserNumber: roomInfo.RoomData.MaxUserNumber,
		GameCount:     roomInfo.RoomData.GameCount,
		UserNumber:    roomInfo.RoomData.UserNumber,
		RoomOwner:     roomInfo.RoomData.RoomOwner,
		RoomWait:      roomInfo.RoomData.RoomWait,
		RoomName:      roomInfo.RoomData.RoomName,
		Users:         roomInfo.RoomData.Users,
	}
	return roomResponse
}

// RoomInfo 房间信息
func (roomInfo *RoomStruct) RoomInfo(message model.Message) {
	//zap.S().Info("[RoomInfo]:收到信息，")
	SendMsgToUser(UsersConn[message.UserID], response.MessageResponse{
		MsgType:  response.RoomInfoResponseType,
		RoomInfo: roomInfo.MakeRoomResponse(),
	})
}

// QuitRoom 退出房间（房主退出会房主转移）
func (roomInfo *RoomStruct) QuitRoom(message model.Message) {
	delete(roomInfo.RoomData.Users, message.UserID)
	roomInfo.RoomData.UserNumber--
	if roomInfo.RoomData.UserNumber == 0 {
		//没人了，销毁房间
		roomInfo.ExitChan <- RoomQuit
		SendMsgToUser(UsersConn[message.UserID], response.MessageResponse{
			MsgType:  response.RoomInfoResponseType,
			RoomInfo: roomInfo.MakeRoomResponse(),
		})
		UsersConn[message.UserID].CloseConn()
		return
	}
	if message.UserID == roomInfo.RoomData.RoomOwner {
		//是房主,转移房间
		num := rand.Intn(int(roomInfo.RoomData.UserNumber))
		for _, data := range roomInfo.RoomData.Users {
			if num <= 0 {
				roomInfo.RoomData.RoomOwner = data.ID
				break
			}
			num--
		}
	}
	BroadcastToAllRoomUsers(roomInfo, response.MessageResponse{
		MsgType:  response.RoomInfoResponseType,
		RoomInfo: roomInfo.MakeRoomResponse(),
	})
	BroadcastToAllRoomUsers(roomInfo, response.MessageResponse{
		MsgType: response.MsgResponseType,
		MsgInfo: &response.MsgResponse{MsgData: "某玩家退出房间"},
	})
	UsersConn[message.UserID].CloseConn()
	//玩家退出，应该从redis删除其服务器连接信息
	global.GameSrvClient.DelConnData(context.Background(), &game_proto.DelConnInfo{
		Id: message.UserID,
	})
}

// UpdateRoom 更新房间的房主或者游戏配置(仅房主)
func (roomInfo *RoomStruct) UpdateRoom(message model.Message) {
	data := message.UpdateData
	if message.UserID != roomInfo.RoomData.RoomOwner {
		//非房主，不可以修改房间的！
		return
	}
	if data.MaxUserNumber >= roomInfo.RoomData.UserNumber && data.MaxUserNumber != 0 {
		roomInfo.RoomData.MaxUserNumber = data.MaxUserNumber
	}
	if data.GameCount != 0 {
		roomInfo.RoomData.GameCount = data.GameCount
	}
	if data.Kicker != 0 {
		//先t人
		if _, ok := roomInfo.RoomData.Users[data.Kicker]; ok {
			//找到人
			delete(roomInfo.RoomData.Users, data.Kicker) //即使找不到人也不报错
			roomInfo.RoomData.UserNumber--
			if UsersConn[data.Kicker] != nil {
				UsersConn[data.Kicker].CloseConn() //可能有nil错误
			}
			global.GameSrvClient.DelConnData(context.Background(), &game_proto.DelConnInfo{
				Id: data.Kicker,
			})
			if roomInfo.RoomData.UserNumber <= 0 {
				roomInfo.ExitChan <- RoomQuit
				//return
			}
		}
	}
	if data.Owner != 0 {
		//查询这个人在不在房间
		if _, ok := roomInfo.RoomData.Users[data.Owner]; ok {
			roomInfo.RoomData.RoomOwner = data.Owner
		}
	}

	SendMsgToUser(UsersConn[message.UserID], response.MessageResponse{
		MsgType: response.MsgResponseType,
		MsgInfo: &response.MsgResponse{
			MsgData: "更新房间成功",
		},
	})
	//更新房间，发送广播
	BroadcastToAllRoomUsers(roomInfo, response.MessageResponse{
		MsgType:  response.RoomInfoResponseType,
		RoomInfo: roomInfo.MakeRoomResponse(),
	})
}

// UpdateUserReadyState 玩家准备状态
func (roomInfo *RoomStruct) UpdateUserReadyState(message model.Message) {
	t := roomInfo.RoomData.Users[message.UserID]
	t.Ready = message.ReadyStateData.IsReady

	SendMsgToUser(UsersConn[message.UserID], response.MessageResponse{
		MsgType: response.MsgResponseType,
		MsgInfo: &response.MsgResponse{
			MsgData: fmt.Sprintf("玩家%d准备状态更新", message.UserID),
		},
	})
	BroadcastToAllRoomUsers(roomInfo, response.MessageResponse{
		MsgType:  response.RoomInfoResponseType,
		RoomInfo: roomInfo.MakeRoomResponse(),
	})
}

// BeginGame 开始游戏
func (roomInfo *RoomStruct) BeginGame(message model.Message) {
	if message.UserID != roomInfo.RoomData.RoomOwner {
		return
	}
	if roomInfo.RoomData.UserNumber != roomInfo.RoomData.MaxUserNumber {
		SendMsgToUser(UsersConn[message.UserID], response.MessageResponse{
			MsgType: response.MsgResponseType,
			MsgInfo: &response.MsgResponse{
				MsgData: "房间没满人,请改房间人数开始游戏",
			},
		})
		return
	}
	for _, data := range roomInfo.RoomData.Users {
		if data.Ready == false {
			SendMsgToUser(UsersConn[message.UserID], response.MessageResponse{
				MsgType: response.MsgResponseType,
				MsgInfo: &response.MsgResponse{
					MsgData: "还有玩家未准备，快T了他",
				},
			})
			return
		}
	}

	user := roomInfo.RoomData.Users[message.UserID]
	user.Ready = true
	roomInfo.RoomData.RoomWait = false
	roomInfo.ExitChan <- GameStart

	BroadcastToAllRoomUsers(roomInfo, response.MessageResponse{
		MsgType: response.MsgResponseType,
		MsgInfo: &response.MsgResponse{
			MsgData: "游戏即将开始",
		},
	})
}

func (roomInfo *RoomStruct) ChatProcess(message model.Message) {
	BroadcastToAllRoomUsers(roomInfo, response.MessageResponse{
		MsgType: response.MsgResponseType,
		MsgInfo: &response.MsgResponse{
			MsgData: fmt.Sprintf("用户%d说：%s", message.UserID, message.ChatMsgData.Data),
		},
	})
}

func (roomInfo *RoomStruct) CheckHealth(message model.Message) {
	//SendMsgToUser(UsersConn[message.UserID], response.MessageResponse{
	//	MsgType:         response.CheckHealthType,
	//	HealthCheckInfo: &response.HealthCheck{},
	//})
}
