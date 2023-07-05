package api

import (
	"context"
	"fmt"
	"game_web/global"
	"game_web/model"
	game_proto "game_web/proto/game"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CHAN 房间号对应创建读取协程的管道
var CHAN = make(map[uint32]chan uint32)

var IntoRoomChan = make(chan bool, 3)

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
			"roomID":        room.RoomID,
			"maxUserNumber": room.MaxUserNumber,
			"gameCount":     room.GameCount,
			"userNumber":    room.UserNumber,
			"roomOwner":     room.RoomOwner,
			"roomWait":      true,
			"roomName":      room.RoomName,
			"users":         user,
		})
		//fmt.Println(user)
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

// 查询房间对应的服务器
func SelectRoomServer(ctx *gin.Context) {
	roomIDStr := ctx.DefaultQuery("room_id", "0")
	roomID, _ := strconv.Atoi(roomIDStr)
	server, err := global.GameSrvClient.GetRoomServer(context.Background(), &game_proto.RoomIDInfo{RoomID: uint32(roomID)})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": fmt.Sprintf("找不到目标服务器%s", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"server": server,
	})
}

// 创建房间，做路由转发
func CreateRoom(ctx *gin.Context) {
	config := api.DefaultConfig()
	config.Address = "139.159.234.134:8500"
	c, _ := api.NewClient(config)
	service, _, err := c.Catalog().Service("process-web", "", nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": "找不到游戏服务器",
		})
		return
	}
	num := rand.Intn(len(service))

	ctx.Redirect(http.StatusPermanentRedirect,
		fmt.Sprintf("http://%s:%d/v1/createRoom", service[num].ServiceAddress, service[num].ServicePort))
}

// 进入房间，做路由转发
func UserIntoRoom(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	//查询该房间对应的路由
	server, err := global.GameSrvClient.GetRoomServer(context.Background(), &game_proto.RoomIDInfo{RoomID: uint32(roomID)})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}
	ctx.Redirect(http.StatusPermanentRedirect,
		fmt.Sprintf("http://%s/v1/userIntoRoom?room_id=%d", server.ServerInfo, roomID))
}
