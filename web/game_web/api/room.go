package api

import (
	"context"
	"encoding/json"
	"errors"
	"game_web/global"
	"game_web/global/response"
	"game_web/model"
	"game_web/proto"
	"game_web/utils"
	"go.uber.org/zap"
)

type DealFunc func(roomID uint32, message model.Message)

var dealFuncs = make(map[uint32]DealFunc)

func init() {
	dealFuncs[model.DeleteRoom] = DropRoom
	dealFuncs[model.GetRoomData] = RoomInfo
	dealFuncs[model.RoomBeginGame] = BeginGame
	dealFuncs[model.UserReadyState] = UpdateUserReadyState
	dealFuncs[model.UpdateRoom] = UpdateRoom
}

// 全局房间信息
var RoomData map[uint32]*model.RoomCoon = make(map[uint32]*model.RoomCoon) //房间号->房间数据的映射(每个房间线程访问各自数据)

func GrpcModelToResponse(room *proto.RoomInfo) response.RoomResponse {
	resp := response.RoomResponse{
		RoomID:        room.RoomID,
		MaxUserNumber: room.MaxUserNumber,
		GameCount:     room.GameCount,
		UserNumber:    room.UserNumber,
		RoomOwner:     room.RoomOwner,
		RoomWait:      room.RoomWait,
	}
	for _, user := range room.Users {
		resp.Users = append(resp.Users, response.UserResponse{
			ID:    user.ID,
			Ready: user.Ready,
		})
	}

	return resp
}

// 房间协程函数(主逻辑)
func startRoomThread(roomID uint32) {
	//初始化房间信息
	if RoomData[roomID] == nil {
		RoomData[roomID] = &model.RoomCoon{
			MsgChan:   make(chan model.Message, 50),
			UsersConn: make(map[uint32]*model.WSConn),
		}
	}
	roomInfo := RoomData[roomID]
	//服务器是单线程处理游戏，那么每次都将客户端发来数据拿过来
	go ReadRoomData(roomInfo)
	//协程主要作用在于处理房间内用户websocket的消息
	for {
		select {
		case msg := <-roomInfo.MsgChan:
			zap.S().Info("收到", msg)
			dealFuncs[msg.Type](roomID, msg)
		case <-roomInfo.PauseChan:
			//停止信号,等待恢复,此时msgchan还在放数据，应该告知发送端不要发了
			<-roomInfo.RecoverChan
			//重新打开房间读取通道,及读取协程
			roomInfo.MsgChan = make(chan model.Message, 50)
			go ReadRoomData(roomInfo)
		}
	}
}

func ReadRoomData(roomInfo *model.RoomCoon) {
	for true {
		for userID, wsConn := range roomInfo.UsersConn {
			if wsConn.IsClose() {
				continue
			}
			data, err := wsConn.InChanRead()
			if err != nil {
				//如果读到客户端关闭信息,关闭与客户端的websocket连接
				wsConn.CloseConn()
				continue
			}
			message := model.Message{}
			err = json.Unmarshal(data, &message)
			if err != nil {
				//客户端发过来数据有误
				zap.S().Info("客户端发送数据有误:", data)
				continue
			}
			message.UserID = userID //添加标识，能够识别用户
			select {
			case <-roomInfo.PauseChan:
				//收到停止信号,停止房间读取线程,关闭传输通道(因为是一对一的通道模式，没其他用户，关闭对其他用户发送没影响)
				close(roomInfo.MsgChan)
				return
			case roomInfo.MsgChan <- message:
				//不断接受信息发送客户端
			}
		}
	}
}

// 房间信息
func RoomInfo(roomID uint32, message model.Message) {
	users := RoomData[roomID].UsersConn
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[UpdateRoom]", err)
		return
	}
	resp := GrpcModelToResponse(room)
	marshal, _ := json.Marshal(resp)
	utils.SendMsgToUser(users[message.UserID], marshal)
}

// 删除房间（仅房主）
func DropRoom(roomID uint32, message model.Message) {
	//先查询房间是否存在
	users := RoomData[roomID].UsersConn
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[UpdateRoom]", err)
		return
	}
	if room.RoomOwner != message.UserID {
		utils.SendErrToUser(users[message.UserID], "[UpdateRoom]", errors.New("非房主不可删除房间"))
		return
	}
	_, err = global.GameSrvClient.DeleteRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[UpdateRoom]", err)
		return
	}
	ret := "删除房间成功"
	err = users[message.UserID].OutChanWrite([]byte(ret))
	if err != nil {
		users[message.UserID].CloseConn()
	}
	//房间变化，广播
	searchRoom, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[UpdateRoom]", err)
		return
	}
	resp := GrpcModelToResponse(searchRoom)
	marshal, _ := json.Marshal(resp)
	BroadcastToAllRoomUsers(RoomData[roomID], marshal)
}

// 更新房间的房主或者游戏配置(仅房主)
func UpdateRoom(roomID uint32, message model.Message) {
	//先查询房间是否存在
	users := RoomData[roomID].UsersConn
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[UpdateRoom]", err)
		return
	}
	if room.RoomOwner != message.UserID {
		utils.SendErrToUser(users[message.UserID], "[UpdateRoom]", errors.New("非房主不可修改房间"))
		return
	}

	roomUpdate := proto.RoomInfo{}
	roomUpdate.RoomID = room.RoomID
	roomUpdate.MaxUserNumber = room.MaxUserNumber
	if message.UpdateData.MaxUserNumber != 0 {
		roomUpdate.MaxUserNumber = message.UpdateData.MaxUserNumber
	}
	roomUpdate.GameCount = room.GameCount
	if message.UpdateData.GameCount != 0 {
		roomUpdate.GameCount = message.UpdateData.GameCount
	}
	roomUpdate.UserNumber = room.UserNumber
	roomUpdate.RoomOwner = room.RoomOwner
	if message.UpdateData.Owner != 0 {
		roomUpdate.RoomOwner = message.UpdateData.Owner
	}
	//T人
	roomUpdate.Users = room.Users
	if message.UpdateData.Kicker != 0 {
		for i, user := range roomUpdate.Users {
			if user.ID == message.UpdateData.Kicker {
				roomUpdate.Users = append(roomUpdate.Users[:i], roomUpdate.Users[i+1:]...)
			}
		}
		roomUpdate.UserNumber--
		//t人了还需要关闭房间里面的连接
		users[message.UpdateData.Kicker].CloseConn()
		delete(users, message.UpdateData.Kicker)
	}

	_, err = global.GameSrvClient.UpdateRoom(context.Background(), &roomUpdate)
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[UpdateRoom]", err)
		return
	}
	err = users[message.UserID].OutChanWrite([]byte("更新房间成功"))
	if err != nil {
		users[message.UserID].CloseConn()
	}
	//更新房间，发送广播
	searchRoom, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[UpdateRoom]", err)
		return
	}
	resp := GrpcModelToResponse(searchRoom)
	marshal, _ := json.Marshal(resp)
	BroadcastToAllRoomUsers(RoomData[roomID], marshal)
}

// 玩家准备状态
func UpdateUserReadyState(roomID uint32, message model.Message) {
	//先查询房间是否存在
	users := RoomData[roomID].UsersConn
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[UpdateUserReadyState]", err)
		return
	}
	//users[message.UserID].IsReady = message.ReadyState.IsReady
	for _, user := range room.Users {
		if user.ID == message.UserID {
			user.Ready = message.ReadyState.IsReady
		}
	}
	utils.SendMsgToUser(users[message.UserID], []byte("玩家准备状态更新成功"))
	//更新房间，发送广播
	searchRoom, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[UpdateUserReadyState]", err)
		return
	}
	resp := GrpcModelToResponse(searchRoom)
	marshal, _ := json.Marshal(resp)
	BroadcastToAllRoomUsers(RoomData[roomID], marshal)
}

// 开始游戏按键接口(需要将用户连接（全局）传输到游戏线程，当前线程就不再运行，等待游戏完成信号再启动该线程，设置pausechan)
func BeginGame(roomID uint32, message model.Message) {
	//查看房间是否存在
	users := RoomData[roomID].UsersConn
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[BeginGame]", err)
		return
	}
	//检查是否是房主
	if room.RoomOwner != message.UserID {
		err := errors.New("非房主不可更新房间")
		utils.SendErrToUser(users[message.UserID], "[BeginGame]", err)
		return
	}
	//检查是否够人了
	if room.UserNumber != room.MaxUserNumber {
		utils.SendMsgToUser(users[message.UserID], []byte("人数不足，无法开始"))
		return
	}
	//都准备好了，可以进入游戏模块,向所有用户发送游戏开始
	BroadcastToAllRoomUsers(RoomData[roomID], []byte("游戏环节开始"))
	//游戏开始,告知房间线程先暂停
	RoomData[roomID].PauseChan <- struct{}{}
	go RunGame(roomID)
}
