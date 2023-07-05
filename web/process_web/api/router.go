package api

import (
	"context"
	"fmt"
	"net/http"
	"process_web/forms"
	"process_web/global"
	"process_web/model"
	"process_web/model/response"
	game_proto "process_web/proto/game"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ConnectCHAN 房间号对应创建读取协程的管道
var ConnectCHAN = make(map[uint32]chan uint32)

// IntoRoomCHAN 用户进房发送chan 房间服务器读取并处理 key:房间号 value:用户id
var IntoRoomCHAN = make(map[uint32]chan uint32)

// IntoRoomRspCHAN IntoRoomChan 对用户进房做出回复  key:房间号 value:加入是否成功
var IntoRoomRspCHAN = make(map[uint32]chan bool)

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

// ConnSocket 建立长连接 TODO 其他非玩家用户进房应该被拒绝
func ConnSocket(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.DefaultQuery("room_id", "0"))
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	if ConnectCHAN[uint32(roomID)] == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "传入room_id错误",
		})
		return
	}
	// 建立websocket连接
	conn, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "无法连接房间服务器",
		})
		return
	}
	if UsersConn[userID] != nil {
		UsersConn[userID].CloseConn()
	}
	UsersConn[userID] = InitWebSocket(conn, userID)
	ConnectCHAN[uint32(roomID)] <- userID
}

// CreateRoom 创建房间,房间创建，需要创建一个协程处理房间及游戏内所有信息 // TODO 创建房间应该先查询房间是否存在
func CreateRoom(ctx *gin.Context) {
	//u1 := ctx.DefaultPostForm("room_id", "默认名字")
	//u2 := ctx.DefaultPostForm("max_user_number", "默认名字")
	//u3 := ctx.DefaultPostForm("game_count", "默认名字")
	//u4 := ctx.DefaultPostForm("room_name", "默认名字")
	//fmt.Println(u1, u2, u3, u4)
	form := forms.CreateRoomForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	zap.S().Infof("[CreateRoom]房间ID:%d", form.RoomID)
	if form.RoomID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "房间号不能小于0",
		})
		return
	}
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID

	var users []*game_proto.RoomUser
	users = append(users, &game_proto.RoomUser{ID: userID, Ready: false})
	zap.S().Infof("[CreateRoom]:注册房间主机和端口%s", fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
	// 1.创建房间对应服务器信息创建成功 //TODO （查询之前是否已经用了该信息）
	_, err := global.GameSrvClient.RecordRoomServer(context.Background(), &game_proto.RecordRoomServerInfo{
		RoomID:     uint32(form.RoomID),
		ServerInfo: fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port),
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}
	// 2.调用创建房间接口
	global.GameSrvClient.SetGlobalRoom(context.Background(), &game_proto.RoomInfo{
		RoomID:        uint32(form.RoomID),
		MaxUserNumber: uint32(form.MaxUserNumber),
		GameCount:     uint32(form.GameCount),
		UserNumber:    1,
		RoomOwner:     userID,
		RoomWait:      true,
		Users:         users,
		RoomName:      form.RoomName,
	})
	// 3.然后创建用户对应服务器的连接
	global.GameSrvClient.RecordConnData(context.Background(), &game_proto.RecordConnInfo{
		ServerInfo: fmt.Sprintf("%s:%d?%d", global.ServerConfig.Host, global.ServerConfig.Port, form.RoomID),
		Id:         userID,
	})
	//启动房间协程
	ConnectCHAN[uint32(form.RoomID)] = make(chan uint32, 10)
	u := make(map[uint32]response.UserData)
	u[userID] = response.UserData{
		ID:    userID,
		Ready: true,
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

// UserIntoRoom 玩家进入房间 房间满人或者其他错误不成功，应该返回错误
func UserIntoRoom(ctx *gin.Context) {
	//zap.S().Infof("[UserIntoRoom]:我在这")
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	//zap.S().Infof("[UserIntoRoom]:RoomID是：%d", roomID)
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	// 玩家进入房间，添加该玩家的服务器连接信息
	zap.S().Infof("[UserIntoRoom]:用户对应的服务器信息%s", fmt.Sprintf("%s:%d?%d", global.ServerConfig.Host, global.ServerConfig.Port, roomID))
	_, err := global.GameSrvClient.RecordConnData(context.Background(), &game_proto.RecordConnInfo{
		ServerInfo: fmt.Sprintf("%s:%d?%d", global.ServerConfig.Host, global.ServerConfig.Port, roomID),
		Id:         userID,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}
	// 告知协程用户进房信息
	zap.S().Infof("[UserIntoRoom]:我进来了%d", uint32(roomID))
	if IntoRoomCHAN[uint32(roomID)] == nil {
		IntoRoomCHAN[uint32(roomID)] = make(chan uint32)
	}
	IntoRoomCHAN[uint32(roomID)] <- userID

	if IntoRoomRspCHAN[uint32(roomID)] == nil {
		IntoRoomRspCHAN[uint32(roomID)] = make(chan bool)
	}
	ok := <-IntoRoomRspCHAN[uint32(roomID)]

	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"err": "进房失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": "ok",
	})
}
