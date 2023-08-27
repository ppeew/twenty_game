package api

import (
	"context"
	"fmt"
	"net/http"
	"process_web/forms"
	"process_web/global"
	"process_web/my_struct"
	"process_web/my_struct/response"
	"process_web/proto/game"
	"process_web/server"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 升级websocket
var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ConnSocket 建立长连接 TODO 其他非玩家用户进房应该被拒绝
func ConnSocket(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.DefaultQuery("room_id", "0"))
	claims, _ := ctx.Get("claims")
	userID := claims.(*my_struct.CustomClaims).ID
	if global.ConnectCHAN[uint32(roomID)] == nil {
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
	value, ok := global.UsersConn.Load(userID)
	if ok {
		value.(*global.WSConn).CloseConn()
	}
	//if global.UsersConn[userID] != nil {
	//	global.UsersConn[userID].CloseConn()
	//}
	global.UsersConn.Store(userID, global.InitWebSocket(conn, userID))
	//global.UsersConn[userID] = global.InitWebSocket(conn, userID)
	global.ConnectCHAN[uint32(roomID)] <- userID
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
	if form.RoomID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "房间号不能小于0",
		})
		return
	}
	claims, _ := ctx.Get("claims")
	userID := claims.(*my_struct.CustomClaims).ID

	var users []*game.RoomUser
	users = append(users, &game.RoomUser{ID: userID, Ready: false})
	//zap.S().Infof("[CreateRoom]:注册房间主机和端口%s", fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
	// 1.创建房间对应服务器信息 //TODO （查询之前是否已经有了该信息，有了就不允许创建）
	_, err := global.GameSrvClient.RecordRoomServer(context.Background(), &game.RecordRoomServerInfo{
		RoomID:     form.RoomID,
		ServerInfo: fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port),
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}
	// 2.然后创建用户对应服务器的连接 TODO (可能因为用户已经再游戏中了创建失败，这个要回滚上一步操作)
	_, err = global.GameSrvClient.RecordConnData(context.Background(), &game.RecordConnInfo{
		ServerInfo: fmt.Sprintf("%s:%d?%d", global.ServerConfig.Host, global.ServerConfig.Port, form.RoomID),
		Id:         userID,
	})
	if err != nil {
		global.GameSrvClient.DelRoomServer(context.Background(), &game.RoomIDInfo{RoomID: uint32(form.RoomID)})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "请先退出原先房间",
		})
		return
	}
	global.ConnectCHAN[form.RoomID] = make(chan uint32, 10)
	u := make(map[uint32]response.UserData)
	u[userID] = response.UserData{
		ID:    userID,
		Ready: true,
	}
	go server.Run(server.NewData(form.RoomID, form.MaxUserNumber, form.GameCount, 1, userID, form.RoomName, []uint32{userID}))
	ctx.JSON(http.StatusOK, gin.H{
		"data": "创建成功",
	})
}

// UserIntoRoom 玩家进入房间 房间满人或者其他错误不成功，应该返回错误
func UserIntoRoom(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	claims, _ := ctx.Get("claims")
	userID := claims.(*my_struct.CustomClaims).ID
	//zap.S().Infof("[UserIntoRoom]:用户对应的服务器信息%s", fmt.Sprintf("%s:%d?%d", global.ServerConfig.Host, global.ServerConfig.Port, roomID))
	// 告知协程用户进房信息
	if global.IntoRoomCHAN[uint32(roomID)] == nil {
		global.IntoRoomCHAN[uint32(roomID)] = make(chan uint32)
	}
	global.IntoRoomCHAN[uint32(roomID)] <- userID

	if global.IntoRoomRspCHAN[uint32(roomID)] == nil {
		global.IntoRoomRspCHAN[uint32(roomID)] = make(chan bool)
	}
	ok := <-global.IntoRoomRspCHAN[uint32(roomID)]

	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"err": "进房失败",
		})
		return
	}
	//记录用户连接信息
	global.GameSrvClient.RecordConnData(context.Background(), &game.RecordConnInfo{
		ServerInfo: fmt.Sprintf("%s:%d?%d", global.ServerConfig.Host, global.ServerConfig.Port, roomID),
		Id:         userID,
	})
	ctx.JSON(http.StatusOK, gin.H{
		"data": "进房成功",
	})
}
