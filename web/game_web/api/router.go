package api

import (
	"context"
	"errors"
	"fmt"
	"game_web/global"
	"game_web/model"
	"game_web/model/response"
	game_proto "game_web/proto/game"
	"game_web/utils"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	NotIn = iota
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

// UsersState 用户ID -> 用户结构体（用户状态+用户连接）
var UsersState = make(map[uint32]*UserState)

type UserState struct {
	State uint32
	WS    *model.WSConn
}

// CreateRoom 创建房间,房间创建，需要创建一个协程处理房间及游戏内所有信息
func CreateRoom(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	roomID, _ := strconv.Atoi(ctx.DefaultQuery("room_id", "0"))
	maxUserNumber, _ := strconv.Atoi(ctx.DefaultQuery("max_user_number", "0"))
	gameCount, _ := strconv.Atoi(ctx.DefaultQuery("game_count", "0"))
	//查询房间
	_, err := global.GameSrvClient.SearchRoom(context.Background(), &game_proto.RoomIDInfo{RoomID: uint32(roomID)})
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"err": "房间已存在，不可创建",
		})
		return
	}
	if UsersState[userID].State != NotIn {
		// 用户在其他房间不能创房
		ctx.JSON(http.StatusOK, gin.H{
			"err": errors.New("请先退出其他房间再创房"),
		})
		return
	}
	_, err = global.GameSrvClient.CreateRoom(context.Background(), &game_proto.RoomInfo{
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
		"data": "创建成功",
	})
	return
}

// UserIntoRoom 玩家进入房间
func UserIntoRoom(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	stateAndConn := UsersState[userID]
	zap.S().Info("用户进入房间，此时状态为：", stateAndConn.State)
	switch stateAndConn.State {
	case RoomIn:
		//用户已经有房间，拒绝进入房间
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "请先退出之前的房间,再进入房间",
		})
	case GameIn:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "正在游戏中，请不要进房",
		})
	case NotIn:
		room, err := global.GameSrvClient.UserIntoRoom(context.Background(), &game_proto.UserIntoRoomInfo{
			RoomID: uint32(roomID),
			UserID: userID,
		})
		if err != nil {
			zap.S().Infof("[UserIntoRoom]:%s", err.Error())
		}
		if room.ErrorMsg != "" {
			ctx.JSON(http.StatusOK, gin.H{
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
		stateAndConn.WS = model.InitWebSocket(conn, userID)
		stateAndConn.State = RoomIn
		// 告知房间主函数，要创建协程来读取用户信息
		CHAN[uint32(roomID)] <- userID
		zap.S().Infof("客户端连接状态:%d", stateAndConn.State)
		// 因为房间更新，给所有订阅者发送房间信息
		BroadcastToAllRoomUsers(room.RoomInfo, GrpcModelToResponse(room.RoomInfo))
		BroadcastToAllRoomUsers(room.RoomInfo, response.RoomMsgResponse{
			MsgType: response.RoomMsgResponseType,
			MsgData: fmt.Sprintf("ID:%d玩家进入房间", userID),
		})
	}
}

// Reconnect 重连游戏服务器
func Reconnect(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	stateAndConn := UsersState[userID]
	if stateAndConn.State == NotIn {
		// 没有需要重连的
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
	stateAndConn.WS = model.InitWebSocket(conn, userID)
	utils.SendMsgToUser(stateAndConn.WS, "重连服务器成功")
}

// SelectUserState 查询用户此时状态-->在房间还是游戏还是没进入房间）
func SelectUserState(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	ctx.JSON(http.StatusOK, gin.H{
		"data": UsersState[userID].State,
	})
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
	})
}
