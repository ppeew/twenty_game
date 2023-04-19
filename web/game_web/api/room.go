package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game_web/global"
	"game_web/global/response"
	"game_web/model"
	"game_web/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

// 获取所有的房间
func GetRoomList(ctx *gin.Context) {
	allRoom, err := global.GameSrvClient.SearchAllRoom(context.Background(), &emptypb.Empty{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var resp []response.RoomResponse
	for _, room := range allRoom.AllRoomInfo {
		r := GrpcModelToResponse(room)
		resp = append(resp, r)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

// 创建房间,房间创建，需要创建一个协程处理房间及游戏内所有信息
func CreateRoom(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID

	roomID, _ := strconv.Atoi(ctx.DefaultQuery("room_id", "0"))
	maxUserNumber, _ := strconv.Atoi(ctx.DefaultQuery("max_user_number", "0"))
	gameCount, _ := strconv.Atoi(ctx.DefaultQuery("game_count", "0"))

	//查询房间
	_, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: uint32(roomID)})
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"err": "房间已存在，不可创建",
		})
		return
	}

	_, err = global.GameSrvClient.CreateRoom(context.Background(), &proto.RoomInfo{
		RoomID:        uint32(roomID),
		MaxUserNumber: uint32(maxUserNumber),
		GameCount:     uint32(gameCount),
		UserNumber:    0,
		RoomOwner:     userID,
		RoomWait:      true,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
		return
	}
	//启动房间协程
	go func(roomID uint32) {
		//初始化房间信息
		publisher := model.NewPublisher()
		roomData := global.RoomData[roomID]
		if global.RoomData[roomID] == nil {
			roomData = new(model.RoomInfo)
			global.RoomData[roomID] = roomData
		}
		roomData.Publisher = publisher

		//协程主要作用在于处理房间内用户websocket的消息
		for {
			select {
			//读信息并处理(msg有类型，比如订阅信息，比如用户消息，用户的更新房间操作)
			case msg := <-roomData.Publisher.MsgChan:
				fmt.Println("收到", msg)
				go HandlerMsg(roomID, msg) //协程处理消息,处理同步问题
			}
		}
		//
	}(uint32(roomID))
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "创建成功",
	})
	return
}

func HandlerMsg(roomID uint32, msg string) {
	data := model.Message{}
	_ = json.Unmarshal([]byte(msg), &data)
	switch data.Type {
	case model.DeleteRoom:
		DropRoom(roomID, data)
	case model.UpdateRoom:
		UpdateRoom(roomID, data)
	case model.GetRoomData:
		RoomInfo(roomID, data)
	case model.UserReadyState:
		UpdateUserReadyState(roomID, data)
	case model.RoomBeginGame:
		BeginGame(roomID, data)
	}
}

// 玩家进入(断线重连)
func UserIntoRoom(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	//查找房间是否存在
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: uint32(roomID)})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
		return
	}

	//断线重连（因为之前已经在房间，查询是否之前有连接过，重连只需要把订阅者内的连接改一下即可）
	for u, _ := range global.RoomData[uint32(roomID)].Publisher.Subscribers {
		if u == userID {
			//存在用户
			conn, err2 := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
			if err2 != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "无法连接房间服务器",
				})
				return
			}
			ws := model.InitWebSocket(conn) //初始化websocket（两协程，分别用来读与写）
			//订阅房间
			sub := model.NewSubscriber(5, ws)
			pub := global.RoomData[uint32(roomID)].Publisher
			pub.AddSubscriber(userID, sub)
			ctx.JSON(http.StatusOK, gin.H{
				"msg": "重连房间服务器成功",
			})
			return
		}
	}

	//房间存在，房间当前人数不应该满了或者房间开始了
	if room.UserNumber >= room.MaxUserNumber {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "房间满了",
		})
		return
	} else if !room.RoomWait {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "房间已开始游戏",
		})
		return
	}

	//进入房间,建立websocket连接
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "无法连接房间服务器",
		})
		return
	}
	//初始化websocket（两协程，分别用来读与写）
	ws := model.InitWebSocket(conn)
	//订阅房间
	sub := model.NewSubscriber(5, ws)
	pub := global.RoomData[uint32(roomID)].Publisher
	pub.AddSubscriber(userID, sub)

	var users []*proto.RoomUser
	users = append(users, &proto.RoomUser{
		ID:    userID,
		Ready: false,
	})
	_, err = global.GameSrvClient.UpdateRoom(context.Background(), &proto.RoomInfo{
		RoomID:        room.RoomID,
		MaxUserNumber: room.MaxUserNumber,
		GameCount:     room.GameCount,
		UserNumber:    room.UserNumber + 1,
		RoomOwner:     room.RoomOwner,
		RoomWait:      room.RoomWait,
		Users:         users,
	})
	if err != nil {
		ws.OutChanWrite([]byte(err.Error()))
	}

	//因为房间更新，给所有订阅者发送房间信息
	searchRoom, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: uint32(roomID)})
	if err != nil {
		ret := err.Error()
		pub.Subscribers[userID].WS.OutChanWrite([]byte(ret))
		return
	}
	resp := GrpcModelToResponse(searchRoom)
	marshal, _ := json.Marshal(resp)
	for _, subscriber := range pub.Subscribers {
		subscriber.WS.OutChanWrite(marshal)
	}
}

// 房间信息
func RoomInfo(roomID uint32, message model.Message) {
	subscribers := global.RoomData[roomID].Publisher.Subscribers
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		ret := err.Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
		return
	}
	resp := GrpcModelToResponse(room)
	marshal, _ := json.Marshal(resp)
	subscribers[message.UserID].WS.OutChanWrite(marshal)
}

// 删除房间（仅房主）
func DropRoom(roomID uint32, message model.Message) {
	//先查询房间是否存在
	subscribers := global.RoomData[roomID].Publisher.Subscribers
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		ret := err.Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
		return
	}
	if room.RoomOwner != message.UserID {
		ret := errors.New("非房主不可删除房间").Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
		return
	}
	_, err = global.GameSrvClient.DeleteRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		ret := err.Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
		return
	}
	ret := "删除房间成功"
	subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
	//房间变化，广播
	searchRoom, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		ret := err.Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
		return
	}
	resp := GrpcModelToResponse(searchRoom)
	marshal, _ := json.Marshal(resp)
	for _, subscriber := range subscribers {
		subscriber.WS.OutChanWrite(marshal)
	}
}

// 更新房间的房主或者游戏配置(仅房主)
func UpdateRoom(roomID uint32, message model.Message) {
	//先查询房间是否存在
	subscribers := global.RoomData[roomID].Publisher.Subscribers
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		ret := err.Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
		return
	}
	if room.RoomOwner != message.UserID {
		ret := errors.New("非房主不可更新房间").Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
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
	}

	_, err = global.GameSrvClient.UpdateRoom(context.Background(), &roomUpdate)
	if err != nil {
		ret := err.Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
		return
	}
	subscribers[message.UserID].WS.OutChanWrite([]byte("更新房间成功"))
	//更新房间，发送广播
	searchRoom, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		ret := err.Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
		return
	}
	resp := GrpcModelToResponse(searchRoom)
	marshal, _ := json.Marshal(resp)
	for _, subscriber := range subscribers {
		subscriber.WS.OutChanWrite(marshal)
	}
}

// 玩家准备状态
func UpdateUserReadyState(roomID uint32, message model.Message) {
	//先查询房间是否存在
	subscribers := global.RoomData[roomID].Publisher.Subscribers
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		ret := err.Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
		return
	}
	//subscribers[message.UserID].IsReady = message.ReadyState.IsReady
	for _, user := range room.Users {
		if user.ID == message.UserID {
			user.Ready = message.ReadyState.IsReady
		}
	}

	subscribers[message.UserID].WS.OutChanWrite([]byte("玩家准备状态更新成功"))
	//更新房间，发送广播
	searchRoom, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		ret := err.Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
		return
	}
	resp := GrpcModelToResponse(searchRoom)
	marshal, _ := json.Marshal(resp)
	for _, subscriber := range subscribers {
		subscriber.WS.OutChanWrite(marshal)
	}
}

// 开始游戏按键接口
func BeginGame(roomID uint32, message model.Message) {
	//查看房间是否存在
	subscribers := global.RoomData[roomID].Publisher.Subscribers
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		ret := err.Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
		return
	}
	//检查是否是房主
	if room.RoomOwner != message.UserID {
		ret := errors.New("非房主不可更新房间").Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
		return
	}
	//检查是否够人了
	if room.UserNumber != room.MaxUserNumber {
		ret := errors.New("人数不足，无法开始").Error()
		subscribers[message.UserID].WS.OutChanWrite([]byte(ret))
		return
	}

	//都准备好了，可以进入游戏模块,发布者向所有用户发送游戏开始，TODO
	for _, subscriber := range subscribers {
		subscriber.WS.OutChanWrite([]byte("游戏环节开始"))
	}
}
