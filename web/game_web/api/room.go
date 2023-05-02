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
	CHAN[roomID] = make(chan uint32, 10)
	return room
}

// 房间主函数
func startRoomThread(roomID uint32) {
	//初始化房间信息
	room := NewRoom(roomID)
	dealFunc := NewDealFunc(room)
	ctx, cancel := context.WithCancel(context.Background())
	searchRoom, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		zap.S().Infof("[startRoomThread]:%s", err)
		cancel()
		return
	}
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
			zap.S().Infof("收到%+v", msg)
			dealFunc[msg.Type](msg)
		case msg := <-room.ExitChan:
			// 停止信号，关闭主函数及读取用户通道函数，优雅退出
			cancel()
			zap.S().Info("[startRoomThread]]:wg正在等待其他协程结束")
			room.wg.Wait()
			zap.S().Info("[startRoomThread]]:其他协程已关闭")
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
		case data := <-UsersState[userID].WS.InChan:
			message := model.Message{}
			err := json.Unmarshal(data, &message)
			if err != nil {
				zap.S().Info("客户端发送数据有误:", string(data))
				utils.SendErrToUser(UsersState[userID].WS, "[ReadRoomData]", err)
				continue
			}
			message.UserID = userID //添加标识，能够识别用户
			roomInfo.MsgChan <- message
		case <-UsersState[userID].WS.CloseChan:
			err := errors.New("连接断开")
			if err != nil {
				//如果与用户的websocket关闭，退出读取协程,并且将该玩家从房间剔除
				roomInfo.wg.Done()
				zap.S().Infof("[ReadRoomUserMsg]:%d用户掉线了", userID)
				roomInfo.MsgChan <- model.Message{
					Type:   model.QuitRoomMsg,
					UserID: userID,
				}
				return
			}
			//default:
			//	data, err := UsersState[userID].WS.InChanRead()
			//	if err != nil {
			//		//如果与用户的websocket关闭，退出读取协程,并且将该玩家从房间剔除
			//		roomInfo.wg.Done()
			//		zap.S().Infof("[ReadRoomUserMsg]:%d用户掉线了", userID)
			//		roomInfo.MsgChan <- model.Message{
			//			Type:   model.QuitRoomMsg,
			//			UserID: userID,
			//		}
			//		return
			//	}
			//	message := model.Message{}
			//	err = json.Unmarshal(data, &message)
			//	if err != nil {
			//		zap.S().Info("客户端发送数据有误:", string(data))
			//		utils.SendErrToUser(UsersState[userID].WS, "[ReadRoomData]", err)
			//		continue
			//	}
			//	message.UserID = userID //添加标识，能够识别用户
			//	roomInfo.MsgChan <- message
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
	utils.SendMsgToUser(UsersState[message.UserID].WS, resp)
}

// QuitRoom 退出房间（房主退出会导致全部房间删除）
func (roomInfo *Room) QuitRoom(message model.Message) {
	//先查询房间是否存在
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomInfo.RoomID})
	if err != nil {
		zap.S().Warnf("[QuitRoom]:%s", err)
		return
	}
	if room.RoomOwner != message.UserID {
		//游戏玩家的退出
		for i, user := range room.Users {
			if message.UserID == user.ID {
				UsersState[message.UserID].WS.CloseConn()
				room.Users = append(room.Users[:i], room.Users[i:]...)
				UsersState[message.UserID].State = NotIn
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
		utils.SendMsgToUser(UsersState[message.UserID].WS, response.RoomMsgResponse{
			MsgType: response.RoomMsgResponseType,
			MsgData: "房主退出房间成功",
		})
		BroadcastToAllRoomUsers(room, "房主退出房间，房间已关闭")
		for _, info := range room.Users {
			UsersState[info.ID].State = NotIn
			UsersState[info.ID].WS.CloseConn()
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
		utils.SendErrToUser(UsersState[message.UserID].WS, "[UpdateRoom]", errors.New("非房主不可修改房间"))
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
			utils.SendErrToUser(UsersState[message.UserID].WS, "[UpdateRoom]", errors.New("房主不可t自己"))
			return
		}
		//发送给被t的玩家
		utils.SendMsgToUser(UsersState[message.UpdateData.Kicker].WS, response.KickerResponse{MsgType: response.KickerResponseType})
		for i, user := range room.Users {
			if user.ID == message.UpdateData.Kicker {
				room.Users = append(room.Users[:i], room.Users[i+1:]...)
			}
		}
		room.UserNumber--
		//t人了还需要关闭房间里面的连接(等待一段时间再关闭连接，为了被t的玩家已经收到被t信息了)
		time.Sleep(1 * time.Second)
		UsersState[message.UpdateData.Kicker].WS.CloseConn()
	}
	_, err = global.GameSrvClient.UpdateRoom(context.Background(), room)
	if err != nil {
		zap.S().Error("[UpdateRoom]:%s", err)
		return
	}
	utils.SendMsgToUser(UsersState[message.UserID].WS, response.RoomMsgResponse{
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
	userInfo := UsersState[message.UserID]
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
	userInfo := UsersState[message.UserID]
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
	room.RoomWait = false
	ownerIndex := uint32(0)
	for i, user := range room.Users {
		if user.Ready == false {
			//没准备好
			if user.ID == room.RoomOwner {
				ownerIndex = uint32(i)
				continue
			}
			utils.SendMsgToUser(userInfo.WS, response.RoomMsgResponse{
				MsgType: response.RoomMsgResponseType,
				MsgData: "其他玩家没准备好",
			})
			return
		}
	}
	room.Users[ownerIndex].Ready = true
	_, err = global.GameSrvClient.UpdateRoom(context.Background(), room)
	for _, info := range room.Users {
		UsersState[info.ID].State = GameIn
	}
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
	utils.SendMsgToUser(UsersState[message.UserID].WS, response.CheckHealthResponse{
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
		err := UsersState[info.ID].WS.OutChanWrite(marshal)
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
