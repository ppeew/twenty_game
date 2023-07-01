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
	"net/http"
	"strconv"
	"strings"

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
var UsersState = make(map[uint32]*WSConn)

// 建立长连接
func ConnSocket(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.DefaultQuery("room_id", "0"))
	if roomID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "传入room_id不能为0",
		})
	}
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	// 建立websocket连接
	conn, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "无法连接房间服务器",
		})
		return
	}
	UsersState[userID] = InitWebSocket(conn, userID)
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
	_, err := global.GameSrvClient.CreateRoom(context.Background(), &game_proto.RoomInfo{
		RoomID:        uint32(form.RoomID),
		MaxUserNumber: uint32(form.MaxUserNumber),
		GameCount:     uint32(form.GameCount),
		UserNumber:    1,
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
	_, err = global.GameSrvClient.RecordConnData(context.Background(), &game_proto.RecordConnInfo{
		ServerInfo: fmt.Sprintf("%s:%d?%d", global.ServerConfig.Host, global.ServerConfig.Port, form.RoomID),
		Id:         userID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	//启动房间协程
	CHAN[uint32(form.RoomID)] = make(chan uint32, 10)
	go startRoomThread(uint32(form.RoomID))
	//CHAN[uint32(form.RoomID)] <- userID
	ctx.JSON(http.StatusOK, gin.H{
		"data": "创建成功",
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
	_, err = global.GameSrvClient.RecordConnData(context.Background(), &game_proto.RecordConnInfo{
		ServerInfo: fmt.Sprintf("%s:%d?%d", global.ServerConfig.Host, global.ServerConfig.Port, roomID),
		Id:         userID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	// 允许进入房间 告知房间主函数，要创建协程来读取用户信息
	//if CHAN[uint32(roomID)] == nil {
	//	CHAN[uint32(roomID)] = make(chan uint32, 10)
	//}
	//CHAN[uint32(roomID)] <- userID
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

// 重连游戏服务器
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
		ctx.JSON(http.StatusBadRequest, gin.H{
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
	UsersState[userID] = InitWebSocket(conn, userID)
	SendMsgToUser(UsersState[userID], response.MessageResponse{
		MsgType: response.MsgResponseType,
		MsgInfo: &response.MsgResponse{MsgData: "重连服务器成功"},
	})
}

// GetRoomInfo 房间信息
func GetRoomInfo(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &game_proto.RoomIDInfo{RoomID: uint32(roomID)})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
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
