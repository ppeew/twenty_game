package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game_web/global"
	"game_web/model"
	"game_web/model/response"
	"game_web/proto"
	"game_web/utils"
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
	return room
}

// 房间主函数
func startRoomThread(roomID uint32) {
	//初始化房间信息
	room := NewRoom(roomID)
	dealFunc := NewDealFunc(room)
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		room.wg.Add(1)
		for true {
			select {
			case userID := <-CHAN:
				go room.ReadRoomUserMsg(ctx, userID)
				room.wg.Add(1)
			case <-ctx.Done():
				room.wg.Done()
				return
			}
		}
	}(ctx)
	//协程主要作用在于处理房间内用户websocket的消息
	for {
		select {
		case msg := <-room.MsgChan:
			zap.S().Infof("收到%+v", msg)
			dealFunc[msg.Type](msg)
		case msg := <-room.ExitChan:
			// 停止信号，关闭主函数及读取用户通道函数，优雅退出
			cancel()
			room.wg.Wait()
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
		default:
			data, err := UsersStateAndConn[userID].WS.InChanRead()
			if err != nil {
				//如果与用户的websocket关闭，退出读取协程
				roomInfo.wg.Done()
				return
			}
			message := model.Message{}
			err = json.Unmarshal(data, &message)
			if err != nil {
				zap.S().Info("客户端发送数据有误:", string(data))
				utils.SendErrToUser(UsersStateAndConn[userID].WS, "[ReadRoomData]", err)
				continue
			}
			message.UserID = userID //添加标识，能够识别用户
			roomInfo.MsgChan <- message
		}
	}
}

// RoomInfo 房间信息
func (roomInfo *Room) RoomInfo(message model.Message) {
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomInfo.RoomID})
	if err != nil {
		zap.S().Error("[RoomInfo]:%s", err)
		return
	}
	resp := GrpcModelToResponse(room)
	utils.SendMsgToUser(UsersStateAndConn[message.UserID].WS, resp)
}

// QuitRoom 退出房间（房主退出会导致全部房间删除）
func (roomInfo *Room) QuitRoom(message model.Message) {
	//先查询房间是否存在
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomInfo.RoomID})
	if err != nil {
		zap.S().Error("[QuitRoom]:%s", err)
		return
	}
	if room.RoomOwner != message.UserID {
		//游戏玩家的退出
		for i, user := range room.Users {
			if message.UserID == user.ID {
				UsersStateAndConn[message.UserID].WS.CloseConn()
				room.Users = append(room.Users[:i], room.Users[i:]...)
				UsersStateAndConn[message.UserID].State = NotIn
				_, err := global.GameSrvClient.UpdateRoom(context.Background(), room)
				if err != nil {
					zap.S().Error("[QuitRoom]错误:%s", err)
				}
				//房间变化，广播
				resp := GrpcModelToResponse(room)
				BroadcastToAllRoomUsers(room, resp)
				break
			}
		}
	} else {
		// 房主退出会销毁房间
		utils.SendMsgToUser(UsersStateAndConn[message.UserID].WS, response.RoomMsgResponse{
			MsgType: response.RoomMsgResponseType,
			MsgData: "房主退出房间成功",
		})
		BroadcastToAllRoomUsers(room, "房主退出房间，房间已关闭")
		for _, info := range room.Users {
			UsersStateAndConn[info.ID].State = NotIn
			UsersStateAndConn[info.ID].WS.CloseConn()
		}
		// 资源释放
		time.Sleep(2 * time.Second)
		_, err = global.GameSrvClient.DeleteRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomInfo.RoomID})
		if err != nil {
			zap.S().Error("[QuitRoom]:%s", err)
			return
		}
		roomInfo.ExitChan <- model.RoomQuit
	}
}

// UpdateRoom 更新房间的房主或者游戏配置(仅房主)
func (roomInfo *Room) UpdateRoom(message model.Message) {
	//先查询房间是否存在
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomInfo.RoomID})
	if err != nil {
		zap.S().Error("[UpdateRoom]:%s", err)
		return
	}
	if room.RoomOwner != message.UserID {
		utils.SendErrToUser(UsersStateAndConn[message.UserID].WS, "[UpdateRoom]", errors.New("非房主不可修改房间"))
		return
	}
	if message.UpdateData.MaxUserNumber != 0 {
		room.MaxUserNumber = message.UpdateData.MaxUserNumber
	}
	if message.UpdateData.GameCount != 0 {
		room.GameCount = message.UpdateData.GameCount
	}
	if message.UpdateData.Owner != 0 {
		room.RoomOwner = message.UpdateData.Owner
	}
	//T人(房主不能t自己)
	if message.UpdateData.Kicker != 0 {
		if message.UpdateData.Kicker == room.RoomOwner {
			utils.SendErrToUser(UsersStateAndConn[message.UserID].WS, "[UpdateRoom]", errors.New("房主不可t自己"))
			return
		}
		//发送给被t的玩家
		utils.SendMsgToUser(UsersStateAndConn[message.UpdateData.Kicker].WS, response.KickerResponse{MsgType: response.KickerResponseType})
		for i, user := range room.Users {
			if user.ID == message.UpdateData.Kicker {
				room.Users = append(room.Users[:i], room.Users[i+1:]...)
			}
		}
		room.UserNumber--
		//t人了还需要关闭房间里面的连接(等待一段时间再关闭连接，为了被t的玩家已经收到被t信息了)
		time.Sleep(1 * time.Second)
		UsersStateAndConn[message.UpdateData.Kicker].WS.CloseConn()
	}
	_, err = global.GameSrvClient.UpdateRoom(context.Background(), room)
	if err != nil {
		zap.S().Error("[UpdateRoom]:%s", err)
		return
	}
	utils.SendMsgToUser(UsersStateAndConn[message.UserID].WS, response.RoomMsgResponse{
		MsgType: response.RoomMsgResponseType,
		MsgData: "更新房间成功",
	})
	//更新房间，发送广播
	resp := GrpcModelToResponse(room)
	BroadcastToAllRoomUsers(room, resp)
}

// UpdateUserReadyState 玩家准备状态
func (roomInfo *Room) UpdateUserReadyState(message model.Message) {
	//先查询房间是否存在
	userInfo := UsersStateAndConn[message.UserID]
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomInfo.RoomID})
	if err != nil {
		zap.S().Error("[UpdateUserReadyState]:%s", err)
		return
	}
	for _, user := range room.Users {
		if user.ID == message.UserID {
			user.Ready = message.ReadyStateData.IsReady
			//更新房间，发送广播
			_, err := global.GameSrvClient.UpdateRoom(context.Background(), room)
			if err != nil {
				zap.S().Error("[UpdateUserReadyState]:%s", err)
				return
			}
			utils.SendMsgToUser(userInfo.WS, response.RoomMsgResponse{
				MsgType: response.RoomMsgResponseType,
				MsgData: fmt.Sprintf("玩家%d准备状态更新", message.UserID),
			})
			resp := GrpcModelToResponse(room)
			BroadcastToAllRoomUsers(room, resp)
			return
		}
	}
	utils.SendErrToUser(userInfo.WS, "[UpdateUserReadyState]", errors.New("没找到该用户"))
}

// BeginGame 开始游戏
func (roomInfo *Room) BeginGame(message model.Message) {
	//查看房间是否存在
	userInfo := UsersStateAndConn[message.UserID]
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomInfo.RoomID})
	if err != nil {
		zap.S().Error("[BeginGame]:%s", err)
		return
	}
	//检查是否是房主
	if room.RoomOwner != message.UserID {
		utils.SendMsgToUser(userInfo.WS, response.RoomMsgResponse{
			MsgType: response.RoomMsgResponseType,
			MsgData: "非房主不可开始游戏",
		})
		return
	}
	//检查是否够人了
	if room.UserNumber != room.MaxUserNumber {
		utils.SendMsgToUser(userInfo.WS, response.RoomMsgResponse{
			MsgType: response.RoomMsgResponseType,
			MsgData: "人数不足，无法开始",
		})
		return
	}
	//游戏开始,房间线程先暂停
	for _, info := range room.Users {
		UsersStateAndConn[info.ID].State = GameIn
	}
	room.RoomWait = false
	_, err = global.GameSrvClient.UpdateRoom(context.Background(), room)
	if err != nil {
		zap.S().Error("[BeginGame]:%s", err)
		return
	}
	roomInfo.ExitChan <- model.GameStart
}

func (roomInfo *Room) ChatProcess(message model.Message) {
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomInfo.RoomID})
	if err != nil {
		zap.S().Errorf("[ChatProcess]:%s", err)
	}
	BroadcastToAllRoomUsers(room, response.RoomMsgResponse{
		MsgType: response.RoomMsgResponseType,
		MsgData: fmt.Sprintf("用户%d说：%s", message.UserID, string(message.ChatMsgData.Data)),
	})
}

func (roomInfo *Room) CheckHealth(message model.Message) {
	utils.SendMsgToUser(UsersStateAndConn[message.UserID].WS, response.CheckHealthResponse{
		MsgType: response.CheckHealthResponseType,
		Ok:      true,
	})
}

func GrpcModelToResponse(roomInfo *proto.RoomInfo) response.RoomResponse {
	resp := response.RoomResponse{
		MsgType:       response.RoomInfoResponseType,
		RoomID:        roomInfo.RoomID,
		MaxUserNumber: roomInfo.MaxUserNumber,
		GameCount:     roomInfo.GameCount,
		UserNumber:    roomInfo.UserNumber,
		RoomOwner:     roomInfo.RoomOwner,
		RoomWait:      roomInfo.RoomWait,
	}
	for _, user := range roomInfo.Users {
		resp.Users = append(resp.Users, response.UserResponse{
			ID:    user.ID,
			Ready: user.Ready,
		})
	}
	return resp
}

func BroadcastToAllRoomUsers(roomInfo *proto.RoomInfo, message interface{}) {
	c := map[string]interface{}{
		"data": message,
	}
	marshal, _ := json.Marshal(c)
	for _, info := range roomInfo.Users {
		err := UsersStateAndConn[info.ID].WS.OutChanWrite(marshal)
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
