package api

import (
	"context"
	"game_web/global"
	"game_web/global/response"
	"game_web/model"
	"game_web/proto"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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
	go startRoomThread(uint32(roomID))

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "创建成功",
	})
	return
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
	for u, _ := range global.RoomData[uint32(roomID)].UsersConn {
		if u == userID {
			//存在用户
			conn, err2 := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
			if err2 != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "无法连接房间服务器",
				})
				return
			}
			//初始化websocket（两协程，分别用来读与写）
			ws := model.InitWebSocket(conn)
			if global.RoomData[uint32(roomID)].UsersConn[userID] == nil {
				global.RoomData[uint32(roomID)].UsersConn[userID] = new(model.WSConn)
			}
			global.RoomData[uint32(roomID)].UsersConn[userID] = ws
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

	if global.RoomData[uint32(roomID)].UsersConn[userID] == nil {
		global.RoomData[uint32(roomID)].UsersConn[userID] = new(model.WSConn)
	}
	global.RoomData[uint32(roomID)].UsersConn[userID] = ws

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
		err := ws.OutChanWrite([]byte(err.Error()))
		if err != nil {
			global.RoomData[uint32(roomID)].UsersConn[userID].CloseConn()
		}
	}
	//因为房间更新，给所有订阅者发送房间信息
	message := model.Message{
		UserID:     userID,
		Type:       model.GetRoomData,
		DeleteData: model.DeleteData{},
		UpdateData: model.UpdateData{},
		RoomData:   model.RoomData{RoomID: uint32(roomID)},
		ReadyState: model.ReadyState{},
		BeginGame:  model.BeginGame{},
	}
	global.RoomData[uint32(roomID)].MsgChan <- message
}
