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
	RoomID    uint32
	Mutex     sync.Mutex               //同时进入房间抢夺房间人数资源
	MsgChan   chan model.Message       //接受信息管道
	ExitChan  chan struct{}            //用于结束房间协程
	ReadExit  chan struct{}            //告知读用户消息线程是否已经完成退出
	UsersConn map[uint32]*model.WSConn //用户id到订阅者的映射
}

// 房间号 -> 房间数据的映射(每个房间线程访问各自数据)
var RoomData map[uint32]*Room = make(map[uint32]*Room)

// 房间主函数(主逻辑)
func startRoomThread(roomID uint32) {
	//初始化房间信息
	roomInfo := NewRoom(roomID)
	//房间对要求实时性不高，采用一个消费者去拿websocket到的数据
	ctx, cancel := context.WithCancel(context.Background())
	go roomInfo.ReadRoomData(ctx)
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
				return
			}
		}
	}
}

func NewRoom(roomID uint32) *Room {
	if RoomData[roomID] == nil {
		RoomData[roomID] = &Room{
			RoomID:    roomID,
			MsgChan:   make(chan model.Message, 1024),
			UsersConn: make(map[uint32]*model.WSConn),
		}
	}
	return RoomData[roomID]
}

// 读取发送到房间的信息入管道
func (roomInfo *Room) ReadRoomData(ctx context.Context) {
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
				zap.S().Info("客户端发送数据有误:", string(data))
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

func init() {
	dealFuncs[model.CheckHealthMsg] = CheckHealth
	dealFuncs[model.QuitRoomMsg] = QuitRoom
	dealFuncs[model.GetRoomMsg] = RoomInfo
	dealFuncs[model.RoomBeginGameMsg] = BeginGame
	dealFuncs[model.UserReadyStateMsg] = UpdateUserReadyState
	dealFuncs[model.UpdateRoomMsg] = UpdateRoom
	dealFuncs[model.ChatMsg] = ChatProcess
}

type dealFunc func(roomID uint32, message model.Message)

var dealFuncs = make(map[uint32]dealFunc)

// 房间信息
func RoomInfo(roomID uint32, message model.Message) {
	users := RoomData[roomID].UsersConn
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[UpdateRoom]", err)
		return
	}
	resp := GrpcModelToResponse(room)
	utils.SendMsgToUser(users[message.UserID], resp)
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
		UsersState[message.UserID] = NotIn
		_, err := global.GameSrvClient.UpdateRoom(context.Background(), room)
		if err != nil {
			zap.S().Infof("[QuitRoom]错误:%s", err)
		}
		//房间变化，广播
		resp := GrpcModelToResponse(room)
		BroadcastToAllRoomUsers(RoomData[roomID], resp)
	} else {
		// 房主退出会销毁房间
		_, err = global.GameSrvClient.DeleteRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
		if err != nil {
			utils.SendErrToUser(users[message.UserID], "[DropRoom]", err)
			return
		}
		utils.SendMsgToUser(users[message.UserID], response.RoomMsgResponse{
			MsgType: response.RoomMsgResponseType,
			MsgData: "房主退出房间成功",
		})
		// 在全局变量内存中删除,防止浪费空间(让房间主线程停下，包括其中的读取队列,还有用户的连接)
		RoomData[roomID].ExitChan <- struct{}{}
		time.Sleep(2 * time.Second)
		<-RoomData[roomID].ReadExit
		for u, conn := range RoomData[roomID].UsersConn {
			UsersState[u] = NotIn
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
		//发送给被t的玩家
		utils.SendMsgToUser(users[message.UpdateData.Kicker], response.KickerResponse{MsgType: response.KickerResponseType})
		for i, user := range room.Users {
			if user.ID == message.UpdateData.Kicker {
				room.Users = append(room.Users[:i], room.Users[i+1:]...)
			}
		}
		room.UserNumber--
		//t人了还需要关闭房间里面的连接(要保证被t的玩家已经收到被t信息了)
		users[message.UpdateData.Kicker].CloseConn()
		delete(users, message.UpdateData.Kicker)
	}

	_, err = global.GameSrvClient.UpdateRoom(context.Background(), room)
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[UpdateRoom]", err)
		return
	}
	utils.SendMsgToUser(users[message.UserID], response.RoomMsgResponse{
		MsgType: response.RoomMsgResponseType,
		MsgData: "更新房间成功",
	})
	//更新房间，发送广播
	resp := GrpcModelToResponse(room)
	BroadcastToAllRoomUsers(RoomData[roomID], resp)
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
		utils.SendMsgToUser(users[message.UserID], response.RoomMsgResponse{
			MsgType: response.RoomMsgResponseType,
			MsgData: fmt.Sprintf("玩家%d准备状态更新", message.UserID),
		})
		resp := GrpcModelToResponse(room)
		BroadcastToAllRoomUsers(RoomData[roomID], resp)
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
		utils.SendMsgToUser(users[message.UserID], response.RoomMsgResponse{
			MsgType: response.RoomMsgResponseType,
			MsgData: "非房主不可开始游戏",
		})
		return
	}
	//检查是否够人了
	if room.UserNumber != room.MaxUserNumber {
		utils.SendMsgToUser(users[message.UserID], response.RoomMsgResponse{
			MsgType: response.RoomMsgResponseType,
			MsgData: "人数不足，无法开始",
		})
		return
	}
	//游戏开始,房间线程先暂停
	for u, _ := range RoomData[roomID].UsersConn {
		UsersState[u] = GameIn
	}
	room.RoomWait = false
	_, err = global.GameSrvClient.UpdateRoom(context.Background(), room)
	if err != nil {
		utils.SendErrToUser(users[message.UserID], "[BeginGame]", err)
		return
	}
	go RunGame(roomID)
	RoomData[roomID].ExitChan <- struct{}{}
}

func ChatProcess(roomID uint32, message model.Message) {
	BroadcastToAllRoomUsers(RoomData[roomID], response.RoomMsgResponse{
		MsgType: response.RoomMsgResponseType,
		MsgData: fmt.Sprintf("用户%d说：%s", message.UserID, string(message.ChatMsgData.Data)),
	})
}

func CheckHealth(roomID uint32, message model.Message) {
	utils.SendMsgToUser(RoomData[roomID].UsersConn[message.UserID], response.CheckHealthResponse{
		MsgType: response.CheckHealthResponseType,
		Ok:      true,
	})
}

func GrpcModelToResponse(room *proto.RoomInfo) response.RoomResponse {
	resp := response.RoomResponse{
		MsgType:       response.RoomInfoResponseType,
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

func BroadcastToAllRoomUsers(coon *Room, message interface{}) {
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
