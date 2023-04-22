package api

import (
	"context"
	"encoding/json"
	"errors"
	"game_web/global"
	"game_web/model"
	"game_web/model/response"
	"game_web/proto"
	"game_web/utils"
	"go.uber.org/zap"
	"time"
)

func init() {
	dealFuncs[model.QuitRoomMsg] = QuitRoom
	dealFuncs[model.GetRoomMsg] = RoomInfo
	dealFuncs[model.RoomBeginGameMsg] = BeginGame
	dealFuncs[model.UserReadyStateMsg] = UpdateUserReadyState
	dealFuncs[model.UpdateRoomMsg] = UpdateRoom
}

type DealFunc func(roomID uint32, message model.Message)

var dealFuncs = make(map[uint32]DealFunc)

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

// 房间主函数(主逻辑)
func startRoomThread(roomID uint32) {
	//初始化房间信息
	if RoomData[roomID] == nil {
		RoomData[roomID] = &model.RoomCoon{
			MsgChan:   make(chan model.Message, 1024),
			UsersConn: make(map[uint32]*model.WSConn),
		}
	}
	roomInfo := RoomData[roomID]
	//房间对要求实时性不高，采用一个消费者去拿websocket到的数据
	ctx, cancel := context.WithCancel(context.Background())
	go ReadRoomData(ctx, roomInfo)
	//协程主要作用在于处理房间内用户websocket的消息
	for {
		select {
		case msg := <-roomInfo.MsgChan:
			zap.S().Info("收到", msg)
			dealFuncs[msg.Type](roomID, msg)
		case <-roomInfo.ExitChan:
			//停止信号，关闭主函数及读取用户通道函数，优雅退出
			cancel()
			select {
			case <-roomInfo.ReadExit:
				//读的子线程已经退出。因为游戏线程还需要房间线程的ws连接，所以不关闭资源
				return
			}
		}
	}
}

func ReadRoomData(ctx context.Context, roomInfo *model.RoomCoon) {
	for true {
		for userID, wsConn := range roomInfo.UsersConn {
			data, err := wsConn.InChanRead()
			if err != nil {
				//如果客户端关闭,不读取了
				continue
			}
			message := model.Message{}
			err = json.Unmarshal(data, &message)
			if err != nil {
				zap.S().Info("客户端发送数据有误:", data)
				utils.SendErrToUser(roomInfo.UsersConn[userID], "[ReadRoomData]", err)
				continue
			}
			message.UserID = userID //添加标识，能够识别用户
			select {
			case <-ctx.Done():
				//收到退出信号,关闭传输通道(因为是一生产者对一消费者的通道模式，没其他生产者，关闭对其他生产者没影响)
				close(roomInfo.MsgChan)
				roomInfo.ReadExit <- struct{}{}
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

// 退出房间（房主退出会导致全部房间删除）
func QuitRoom(roomID uint32, message model.Message) {
	//先查询房间是否存在
	users := RoomData[roomID].UsersConn
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[DropRoom]", err)
		return
	}
	if room.RoomOwner != message.UserID {
		//游戏玩家的退出
		for i, user := range RoomData[roomID].UsersConn {
			if i == message.UserID {
				user.CloseConn()
				break
			}
		}
		for i, user := range room.Users {
			if message.UserID == user.ID {
				room.Users = append(room.Users[:i], room.Users[i:]...)
				break
			}
		}
		_, err := global.GameSrvClient.UpdateRoom(context.Background(), room)
		if err != nil {
			zap.S().Infof("[QuitRoom]错误:%s", err)
		}
	} else {
		_, err = global.GameSrvClient.DeleteRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
		if err != nil {
			utils.SendErrToUser(users[message.UserID], "[DropRoom]", err)
			return
		}
		utils.SendMsgToUser(users[message.UserID], "删除房间成功")
		//房间变化，广播
		searchRoom, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
		if err != nil {
			utils.SendErrToUser(users[message.UserID], "[DropRoom]", err)
			return
		}
		resp := GrpcModelToResponse(searchRoom)
		marshal, _ := json.Marshal(resp)
		BroadcastToAllRoomUsers(RoomData[roomID], marshal)
		// 在全局变量内存中删除,防止浪费空间(让房间主线程停下，包括其中的读取队列,还有用户的连接)
		RoomData[roomID].ExitChan <- struct{}{}
		time.Sleep(2 * time.Second)
		<-RoomData[roomID].ReadExit
		for _, conn := range RoomData[roomID].UsersConn {
			conn.CloseConn()
		}
		delete(RoomData, roomID)
	}
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
			utils.SendErrToUser(users[message.UserID], "[UpdateRoom]", errors.New("房主不可t自己"))
			return
		}
		for i, user := range room.Users {
			if user.ID == message.UpdateData.Kicker {
				room.Users = append(room.Users[:i], room.Users[i+1:]...)
			}
		}
		room.UserNumber--
		//t人了还需要关闭房间里面的连接
		users[message.UpdateData.Kicker].CloseConn()
		delete(users, message.UpdateData.Kicker)
	}

	_, err = global.GameSrvClient.UpdateRoom(context.Background(), room)
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[UpdateRoom]", err)
		return
	}
	utils.SendMsgToUser(users[message.UserID], "更新房间成功")
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
	isFind := false
	for _, user := range room.Users {
		if user.ID == message.UserID {
			user.Ready = message.ReadyStateData.IsReady
			isFind = true
		}
	}
	if isFind {
		//更新房间，发送广播
		_, err := global.GameSrvClient.UpdateRoom(context.Background(), room)
		if err != nil {
			utils.SendErrToUser(users[message.UserID], "[UpdateUserReadyState]", err)
			return
		}
		searchRoom, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
		if err != nil {
			utils.SendErrToUser(users[message.UserID], "[UpdateUserReadyState]", err)
			return
		}
		utils.SendMsgToUser(users[message.UserID], "玩家准备状态更新成功")
		resp := GrpcModelToResponse(searchRoom)
		marshal, _ := json.Marshal(resp)
		BroadcastToAllRoomUsers(RoomData[roomID], marshal)
	} else {
		utils.SendErrToUser(users[message.UserID], "[UpdateUserReadyState]", errors.New("没找到该用户"))
	}
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
		utils.SendMsgToUser(users[message.UserID], "非房主不可更新房间")
		return
	}
	//检查是否够人了
	if room.UserNumber != room.MaxUserNumber {
		utils.SendMsgToUser(users[message.UserID], "人数不足，无法开始")
		return
	}
	//都准备好了，可以进入游戏模块,向所有用户发送游戏开始
	BroadcastToAllRoomUsers(RoomData[roomID], []byte("游戏环节开始"))
	//游戏开始,告知房间线程先暂停
	room.RoomWait = false
	_, err = global.GameSrvClient.UpdateRoom(context.Background(), room)
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[BeginGame]", err)
		return
	}
	go RunGame(roomID)
	RoomData[roomID].ExitChan <- struct{}{}
}

func BroadcastToAllRoomUsers(coon *model.RoomCoon, message interface{}) {
	c := map[string]interface{}{
		"data": message,
	}
	marshal, _ := json.Marshal(c)
	for userID, wsConn := range coon.UsersConn {
		err := wsConn.OutChanWrite(marshal)
		if err != nil {
			zap.S().Infof("ID为%d的用户掉线了", userID)
		}
	}
}
