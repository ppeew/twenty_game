package api

import (
	"context"
	"fmt"
	"game_web/forms"
	"game_web/global"
	"game_web/model"
	"game_web/model/response"
	game_proto "game_web/proto/game"
	user_proto "game_web/proto/user"
	"game_web/utils"
	"net/http"
	"strconv"

	"google.golang.org/grpc"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 0->大厅 1->房间 2->游戏
const (
	OutSide = iota
	RoomIn
	GameIn
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

// UsersState 用户ID -> 用户连接
var UsersState = make(map[uint32]*model.WSConn)

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
	_, err := global.GameSrvClient.CreateRoom(context.Background(), &game_proto.RoomInfo{
		RoomID:        uint32(form.RoomID),
		MaxUserNumber: uint32(form.MaxUserNumber),
		GameCount:     uint32(form.GameCount),
		UserNumber:    0,
		RoomOwner:     userID,
		RoomWait:      true,
		RoomName:      form.RoomName,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	//启动房间协程
	go startRoomThread(uint32(form.RoomID))
	ctx.JSON(http.StatusOK, gin.H{
		"data": "创建成功",
		"err":  "",
	})
}

// UserIntoRoom 玩家进入房间
func UserIntoRoom(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	room, err := global.GameSrvClient.UserIntoRoom(context.Background(), &game_proto.UserIntoRoomInfo{
		RoomID: uint32(roomID),
		UserID: userID,
	})
	if err != nil {
		zap.S().Infof("[UserIntoRoom]:%s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}
	if room.ErrorMsg != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": room.ErrorMsg,
		})
		return
	}
	// 允许进入房间,建立websocket连接
	conn, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "无法连接房间服务器",
		})
		return
	}
	UsersState[userID] = model.InitWebSocket(conn, userID)
	//_, err = global.UserSrvClient.UpdateUserState(context.Background(), &user_proto.UpdateUserStateInfo{Id: userID, State: RoomIn})
	if err != nil {
		zap.S().Warnf("[UserIntoRoom]:%s", err)
		return
	}
	// 告知房间主函数，要创建协程来读取用户信息
	CHAN[uint32(roomID)] <- userID
	// 因为房间更新，给所有订阅者发送房间信息
	BroadcastToAllRoomUsers(room.RoomInfo, GrpcModelToResponse(room.RoomInfo))
	BroadcastToAllRoomUsers(room.RoomInfo, response.RoomMsgResponse{
		MsgType: response.RoomMsgResponseType,
		MsgData: fmt.Sprintf("ID:%d玩家进入房间", userID),
	})

}

// 获得重连服务器信息
func GetReconnInfo(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	info, err := global.GameSrvClient.GetReconnInfo(context.Background(), &game_proto.UserIDInfo{Id: userID})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"serverInfo": info.ServerInfo,
	})
}

// Reconnect 重连游戏服务器
func Reconnect(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}
	userSrvClient := user_proto.NewUserClient(userConn)
	state, err := userSrvClient.GetUserState(context.Background(), &user_proto.UserIDInfo{Id: userID})
	if err != nil {
		zap.S().Warnf("[CreateRoom]:%s", err)
		return
	}
	if state.State == OutSide {
		ctx.JSON(http.StatusOK, gin.H{
			"data": "不需要重连",
		})
		return
	}
	conn, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "无法连接房间服务器",
		})
		return
	}
	UsersState[userID] = model.InitWebSocket(conn, userID)
	utils.SendMsgToUser(UsersState[userID], "重连服务器成功")
}

// GetRoomInfo 房间信息
func GetRoomInfo(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &game_proto.RoomIDInfo{RoomID: uint32(roomID)})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"err": err,
		})
		return
	}
	resp := GrpcModelToResponse(room)
	ctx.JSON(http.StatusOK, gin.H{
		"data": resp,
		"err":  "",
	})
}

// GetRoomList 获取所有的房间
func GetRoomList(ctx *gin.Context) {
	allRoom, err := global.GameSrvClient.SearchAllRoom(context.Background(), &emptypb.Empty{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	var resp []response.RoomResponse
	for _, room := range allRoom.AllRoomInfo {
		r := GrpcModelToResponse(room)
		resp = append(resp, r)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": resp,
		"err":  "",
	})
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
		"err":  "",
	})
}
