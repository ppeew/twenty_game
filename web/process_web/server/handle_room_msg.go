package server

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"process_web/global"
	"process_web/my_struct"
	"process_web/my_struct/response"
	"process_web/proto/game"
	"process_web/utils"
	"sort"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"

	"go.uber.org/zap"
)

type dealFunc func(message my_struct.Message)

func NewDealFunc(room *RoomStruct) map[uint32]dealFunc {
	var dealFun = make(map[uint32]dealFunc)
	dealFun[my_struct.CheckHealthMsg] = room.CheckHealth
	dealFun[my_struct.QuitRoomMsg] = room.QuitRoom
	dealFun[my_struct.GetRoomMsg] = room.RoomInfo
	dealFun[my_struct.RoomBeginGameMsg] = room.BeginGame
	dealFun[my_struct.UserReadyStateMsg] = room.UpdateUserReadyState
	dealFun[my_struct.UpdateRoomMsg] = room.UpdateRoom
	dealFun[my_struct.ChatMsg] = room.ChatProcess
	dealFun[my_struct.UserIntoMsg] = room.UserInto //仅服务器用
	return dealFun
}

// RoomInfo 房间信息
func (room *RoomStruct) RoomInfo(message my_struct.Message) {
	global.SendMsgToUser(global.UsersConn[message.UserID], response.MessageResponse{
		MsgType:  response.RoomInfoResponseType,
		RoomInfo: room.MakeRoomResponse(),
	})
}

// QuitRoom 退出房间（房主退出会房主转移）
func (room *RoomStruct) QuitRoom(message my_struct.Message) {
	delete(room.Users, message.UserID)
	room.UserNumber--
	if room.UserNumber == 0 {
		//没人了，销毁房间
		room.ExitChan <- RoomQuit
		global.UsersConn[message.UserID].CloseConn()
		return
	}
	if message.UserID == room.RoomOwner {
		//是房主,转移房间
		num := rand.Intn(int(room.UserNumber))
		for _, data := range room.Users {
			if num <= 0 {
				room.RoomOwner = data.ID
				global.SendMsgToUser(global.UsersConn[data.ID], response.MessageResponse{
					MsgType: response.MsgResponseType,
					MsgInfo: &response.MsgResponse{MsgData: "房主是你的了"},
				})
				break
			}
			num--
		}
	}
	BroadcastToAllRoomUsers(room, response.MessageResponse{
		MsgType:  response.RoomInfoResponseType,
		RoomInfo: room.MakeRoomResponse(),
	})
	global.UsersConn[message.UserID].CloseConn()
	//玩家退出，应该从redis删除其服务器连接信息
	global.GameSrvClient.DelConnData(context.Background(), &game.DelConnInfo{
		Id: message.UserID,
	})
}

// UpdateRoom 更新房间的房主或者游戏配置(仅房主)
func (room *RoomStruct) UpdateRoom(message my_struct.Message) {
	data := message.UpdateData
	if message.UserID != room.RoomOwner {
		//非房主，不可以修改房间的！
		return
	}
	if data.RoomName != "" {
		room.RoomName = data.RoomName
	}
	if data.MaxUserNumber >= room.UserNumber && data.MaxUserNumber != 0 {
		room.MaxUserNumber = data.MaxUserNumber
	}
	if data.GameCount != 0 {
		room.GameCount = data.GameCount
	}
	if data.Kicker != 0 {
		//先t人
		if _, ok := room.Users[data.Kicker]; ok {
			//找到人
			global.SendMsgToUser(global.UsersConn[data.Kicker], response.MessageResponse{
				MsgType: response.KickerResponseType,
				KickerInfo: &response.KickerResponse{
					ID: data.Kicker,
				},
			})
			delete(room.Users, data.Kicker)
			room.UserNumber--
			//if global.UsersConn[data.Kicker] != nil {
			//	global.UsersConn[data.Kicker].CloseConn() //可能有nil错误
			//}
			global.GameSrvClient.DelConnData(context.Background(), &game.DelConnInfo{
				Id: data.Kicker,
			})
			if room.UserNumber <= 0 {
				room.ExitChan <- RoomQuit
			}
		}
	}
	if data.Owner != 0 {
		//查询这个人在不在房间
		if _, ok := room.Users[data.Owner]; ok {
			room.RoomOwner = data.Owner
		}
	}

	global.SendMsgToUser(global.UsersConn[message.UserID], response.MessageResponse{
		MsgType: response.MsgResponseType,
		MsgInfo: &response.MsgResponse{
			MsgData: "更新房间成功",
		},
	})
	//更新房间，发送广播
	BroadcastToAllRoomUsers(room, response.MessageResponse{
		MsgType:  response.RoomInfoResponseType,
		RoomInfo: room.MakeRoomResponse(),
	})
}

// UpdateUserReadyState 玩家准备状态
func (room *RoomStruct) UpdateUserReadyState(message my_struct.Message) {
	t := room.Users[message.UserID]
	t.Ready = message.ReadyStateData.IsReady
	room.Users[message.UserID] = t
	BroadcastToAllRoomUsers(room, response.MessageResponse{
		MsgType:  response.RoomInfoResponseType,
		RoomInfo: room.MakeRoomResponse(),
	})
}

// BeginGame 开始游戏
func (room *RoomStruct) BeginGame(message my_struct.Message) {
	if message.UserID != room.RoomOwner {
		global.SendMsgToUser(global.UsersConn[message.UserID], response.MessageResponse{
			MsgType: response.MsgResponseType,
			MsgInfo: &response.MsgResponse{
				MsgData: fmt.Sprintf("玩家%d不是房主，不能开始游戏", message.UserID),
			},
		})
		return
	}
	if room.UserNumber != room.MaxUserNumber {
		global.SendMsgToUser(global.UsersConn[message.UserID], response.MessageResponse{
			MsgType: response.MsgResponseType,
			MsgInfo: &response.MsgResponse{
				MsgData: "房间没满人,请改房间人数开始游戏",
			},
		})
		return
	}
	for _, data := range room.Users {
		if data.Ready == false && data.ID != room.RoomOwner {
			global.SendMsgToUser(global.UsersConn[message.UserID], response.MessageResponse{
				MsgType: response.MsgResponseType,
				MsgInfo: &response.MsgResponse{
					MsgData: "还有玩家未准备，快T了他",
				},
			})
			return
		}
	}

	user := room.Users[message.UserID]
	user.Ready = true
	room.RoomWait = false
	room.ExitChan <- GameStart

	BroadcastToAllRoomUsers(room, response.MessageResponse{
		MsgType:       response.BeginGameResponseType,
		BeginGameInfo: &response.BeginGameData{},
	})
}

func (room *RoomStruct) ChatProcess(message my_struct.Message) {
	zap.S().Infof("[ChatProcess]:%d,%s", message.UserID, message.ChatMsgData.Data)
	BroadcastToAllRoomUsers(room, response.MessageResponse{
		MsgType: response.ChatResponseType,
		ChatInfo: &response.ChatResponse{
			UserID:      message.UserID,
			ChatMsgData: message.ChatMsgData.Data,
		},
	})
}

func (room *RoomStruct) CheckHealth(message my_struct.Message) {
	global.SendMsgToUser(global.UsersConn[message.UserID], response.MessageResponse{
		MsgType:         response.CheckHealthType,
		HealthCheckInfo: &response.HealthCheck{},
	})
}

// 仅服务器使用的
func (room *RoomStruct) UserInto(message my_struct.Message) {
	success := false
	if _, exist := room.Users[message.UserID]; !exist && room.UserNumber < room.MaxUserNumber {
		room.UserNumber++
		//查询API用户信息
		var res utils.UserInfo
		gorequest.New().Get("http://139.159.234.134:8000/user/v1/search").Param("id", strconv.Itoa(int(message.UserID))).
			Retry(5, time.Second, http.StatusInternalServerError).EndStruct(&res)
		room.Users[message.UserID] = my_struct.UserRoomData{
			ID:           message.UserID,
			Ready:        false,
			IntoRoomTime: time.Now(),
			Nickname:     res.Nickname,
			Gender:       res.Gender,
			Username:     res.Username,
			Image:        res.Image,
		}
		zap.S().Infof("[UserInto]:查询出用户信息%+v", res)
		success = true
	}
	if global.IntoRoomRspCHAN[room.RoomID] == nil {
		global.IntoRoomRspCHAN[room.RoomID] = make(chan bool)
	}
	global.IntoRoomRspCHAN[room.RoomID] <- success
	BroadcastToAllRoomUsers(room, response.MessageResponse{
		MsgType:  response.RoomInfoResponseType,
		RoomInfo: room.MakeRoomResponse(),
	})
}

func (room *RoomStruct) MakeRoomResponse() *response.RoomResponse {
	var users []response.UserData
	for _, data := range room.Users {
		users = append(users, response.UserData{
			ID:           data.ID,
			Ready:        data.Ready,
			IntoRoomTime: data.IntoRoomTime,
			Nickname:     data.Nickname,
			Gender:       data.Gender,
			Username:     data.Username,
			Image:        data.Image,
		})
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].IntoRoomTime.Before(users[j].IntoRoomTime)
	})
	roomResponse := &response.RoomResponse{
		RoomID:        room.RoomID,
		MaxUserNumber: room.MaxUserNumber,
		GameCount:     room.GameCount,
		UserNumber:    room.UserNumber,
		RoomOwner:     room.RoomOwner,
		RoomWait:      room.RoomWait,
		RoomName:      room.RoomName,
		Users:         users,
	}
	return roomResponse
}
