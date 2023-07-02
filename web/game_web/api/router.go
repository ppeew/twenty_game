package api

import (
	"context"
	"fmt"
	"game_web/forms"
	"game_web/global"
	"game_web/model"
	"game_web/model/response"
	game_proto "game_web/proto/game"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CHAN 房间号对应创建读取协程的管道
var CHAN = make(map[uint32]chan uint32)

// 升级websocket
var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// UsersConn 用户ID -> 用户连接
var UsersConn = make(map[uint32]*WSConn)

// 建立长连接
func ConnSocket(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.DefaultQuery("room_id", "0"))
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	if CHAN[uint32(roomID)] == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "传入room_id错误",
		})
	}
	// 建立websocket连接
	conn, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "无法连接房间服务器",
		})
		return
	}
	UsersConn[userID] = InitWebSocket(conn, userID)
	CHAN[uint32(roomID)] <- userID
}

// CreateRoom 创建房间,房间创建，需要创建一个协程处理房间及游戏内所有信息
func CreateRoom(ctx *gin.Context) {
	form := forms.CreateRoomForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID

	var users []*game_proto.RoomUser
	users = append(users, &game_proto.RoomUser{ID: userID, Ready: false})
	_, err := global.GameSrvClient.SetGlobalRoom(context.Background(), &game_proto.RoomInfo{
		RoomID:        uint32(form.RoomID),
		MaxUserNumber: uint32(form.MaxUserNumber),
		GameCount:     uint32(form.GameCount),
		UserNumber:    1,
		RoomOwner:     userID,
		RoomWait:      true,
		Users:         users,
		RoomName:      form.RoomName,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "房间存在",
		})
		return
	}
	global.GameSrvClient.RecordConnData(context.Background(), &game_proto.RecordConnInfo{
		ServerInfo: fmt.Sprintf("%s:%d?%d", global.ServerConfig.Host, global.ServerConfig.Port, form.RoomID),
		Id:         userID,
	})
	//启动房间协程
	CHAN[uint32(form.RoomID)] = make(chan uint32, 10)
	u := make(map[uint32]response.UserData)
	u[userID] = response.UserData{
		ID:    userID,
		Ready: false,
	}
	go startRoomThread(RoomData{
		RoomID:        uint32(form.RoomID),
		MaxUserNumber: uint32(form.MaxUserNumber),
		GameCount:     uint32(form.GameCount),
		UserNumber:    1,
		RoomOwner:     userID,
		RoomWait:      true,
		Users:         u,
		RoomName:      form.RoomName,
	})
	ctx.JSON(http.StatusOK, gin.H{
		"data": "创建成功",
	})
}

// UserIntoRoom 玩家进入房间
func UserIntoRoom(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	global.GameSrvClient.RecordConnData(context.Background(), &game_proto.RecordConnInfo{
		ServerInfo: fmt.Sprintf("%s:%d?%d", global.ServerConfig.Host, global.ServerConfig.Port, roomID),
		Id:         userID,
	})
	//告知协程
	CHAN[uint32(roomID)] <- userID
	ctx.JSON(http.StatusOK, gin.H{
		"data": "ok",
	})
}

// 获得重连服务器信息
func GetConnInfo(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	info, err := global.GameSrvClient.GetConnData(context.Background(), &game_proto.UserIDInfo{Id: userID})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	split := strings.Split(info.ServerInfo, "?")
	ctx.JSON(http.StatusOK, gin.H{
		"serverInfo": split[0],
		"roomID":     split[1],
	})
}

// GetRoomInfo 房间信息
//func GetRoomInfo(ctx *gin.Context) {
//	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
//	room, err := global.GameSrvClient.SearchRoom(context.Background(), &game_proto.RoomIDInfo{RoomID: uint32(roomID)})
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{
//			"err": err,
//		})
//		return
//	}
//	resp := GrpcModelToResponse(room)
//	ctx.JSON(http.StatusOK, gin.H{
//		"data": resp,
//	})
//}

// GetRoomList 获取所有的房间
func GetRoomList(ctx *gin.Context) {
	allRoom, err := global.GameSrvClient.SearchAllRoom(context.Background(), &emptypb.Empty{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	var resp []map[string]interface{}
	for _, room := range allRoom.AllRoomInfo {
		var user []uint32
		for _, roomUser := range room.Users {
			user = append(user, roomUser.ID)
		}
		resp = append(resp, map[string]interface{}{
			"RoomID":        room.RoomID,
			"MaxUserNumber": room.MaxUserNumber,
			"GameCount":     room.GameCount,
			"UserNumber":    room.UserNumber,
			"RoomOwner":     room.RoomOwner,
			"RoomWait":      true,
			"RoomName":      room.RoomName,
			"Users":         user,
		})
		//fmt.Println(user)
	}
	ctx.JSON(http.StatusOK, resp)
}

// SelectItems 查询个人的物品信息
func SelectItems(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	info, err := global.GameSrvClient.GetUserItemsInfo(context.Background(), &game_proto.UserIDInfo{Id: userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": info,
	})
}
