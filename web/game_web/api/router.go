package api

import (
	"context"
	"errors"
	"fmt"
	"game_web/global"
	"game_web/model"
	"game_web/model/response"
	"game_web/proto"
	"game_web/utils"
	"net/http"
	"strconv"
	"sync"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/known/emptypb"
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
		"data": "创建成功",
	})
	return
}

// 玩家进入房间(断线重连)
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
	if RoomData[uint32(roomID)] == nil {
		//没有创房间
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "还没创房就进入",
		})
		return
	}
	//一个用户同时间只能够在一间房（房间或者游戏）存在
	state := UsersState[userID]
	if state == nil {
		state = &UserState{
			State:   NotIn,
			RWMutex: sync.RWMutex{},
		}
	}
	state.RWMutex.RLock()
	if state.State == RoomIn {
		//用户进入相同房间会重连,否则要求先退出该房间
		err = ReConnRoom(RoomData[uint32(roomID)].UsersConn, ctx, userID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"err": "请先退出之前的房间,再创房",
			})
		}
		return
	} else if state.State == GameIn {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "正在游戏中，请不要创房",
		})
		return
	}
	state.RWMutex.RUnlock()
	//房间存在，房间当前人数不应该满了或者房间开始了
	if room.UserNumber >= room.MaxUserNumber {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "房间满了",
		})
		return
	} else if !room.RoomWait {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "房间已开始游戏",
		})
		return
	}
	//进入房间,建立websocket连接
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "无法连接房间服务器",
		})
		return
	}
	//初始化websocket（两协程，分别用来读与写）
	ws := model.InitWebSocket(conn, userID)
	if RoomData[uint32(roomID)].UsersConn[userID] == nil {
		RoomData[uint32(roomID)].UsersConn[userID] = new(model.WSConn)
	}
	RoomData[uint32(roomID)].UsersConn[userID] = ws
	room.Users = append(room.Users, &proto.RoomUser{
		ID:    userID,
		Ready: false,
	})
	room.UserNumber += 1
	_, err = global.GameSrvClient.UpdateRoom(context.Background(), room)
	if err != nil {
		utils.SendErrToUser(ws, "[UserIntoRoom]", err)
	}
	state.State = RoomIn
	//因为房间更新，给所有订阅者发送房间信息
	room, err = global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: uint32(roomID)})
	if err != nil {
		utils.SendErrToUser(ws, "[UserIntoRoom]", err)
	}
	BroadcastToAllRoomUsers(RoomData[uint32(roomID)], GrpcModelToResponse(room))
	BroadcastToAllRoomUsers(RoomData[uint32(roomID)], response.RoomMsgResponse{
		MsgType: response.RoomMsgResponseType,
		MsgData: fmt.Sprintf("ID:%d玩家进入房间", userID),
	})
}

func ReConnRoom(usersConn map[uint32]*model.WSConn, ctx *gin.Context, userID uint32) error {
	//断线重连机制
	for u, _ := range usersConn {
		if u == userID {
			//存在用户,先把之前开的协程关闭
			usersConn[u].CloseConn()
			conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
			if err != nil {
				zap.S().Panic("[ReConnRoom]重连服务器失败")
				return nil
			}
			//初始化websocket（两协程，分别用来读与写）
			ws := model.InitWebSocket(conn, userID)
			if usersConn[userID] == nil {
				usersConn[userID] = new(model.WSConn)
			}
			usersConn[userID] = ws
			utils.SendMsgToUser(ws, "重连房间服务器成功")
			return nil
		}
	}
	return errors.New("该游戏玩家不在房间中")
}

// 房间信息
func GetRoomInfo(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: uint32(roomID)})
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

// 查询个人的物品信息
func SelectItems(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	info, err := global.GameSrvClient.GetUserItemsInfo(context.Background(), &proto.UserIDInfo{Id: userID})
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

// 查询用户此时状态（用于断线重连）（是在房间还是游戏还是没进入房间）
func SelectUserState(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	state := UsersState[userID]
	if state == nil {
		state = &UserState{
			State:   NotIn,
			RWMutex: sync.RWMutex{},
		}
	}
	state.RWMutex.RLock()
	ctx.JSON(http.StatusOK, gin.H{
		"data": state.State,
	})
	state.RWMutex.RUnlock()
}

// 玩家进入游戏(断线重连)
func UserIntoGame(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	if GameData[uint32(roomID)] == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"err": "不存在该游戏房间，稍后再试",
		})
		return
	}
	//等待服务器初始化完成，因为有些资源还没分配好,一旦InitChan读到，说明服务器已经做好了初始化准备
	<-GameData[uint32(roomID)].InitChan
	isFindUser := false
	for u, _ := range GameData[uint32(roomID)].Users {
		if u == userID {
			//找到用户，建立连接
			isFindUser = true
			break
		}
	}
	if !isFindUser {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "该玩家不在该游戏房间",
		})
		return
	}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "无法连接游戏服务器",
		})
		return
	}
	ws := model.InitWebSocket(conn, userID)
	GameData[uint32(roomID)].Users[userID].WS = ws
	utils.SendMsgToUser(ws, "连接游戏服务器成功")
}
