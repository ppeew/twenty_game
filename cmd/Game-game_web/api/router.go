package api

import (
	"context"
	"fmt"
	"game_web/forms"
	"game_web/global"
	"game_web/model"
	game_proto "game_web/proto/game"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/gin-gonic/gin"
)

// CHAN 房间号对应创建读取协程的管道
var CHAN = make(map[uint32]chan uint32)

// 获得重连服务器信息
func GetConnInfo(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	info, err := global.GameSrvClient.GetConnData(context.Background(), &game_proto.UserIDInfo{Id: userID})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"serverInfo": "",
			"roomID":     "",
		})
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
	index, _ := strconv.Atoi(ctx.DefaultQuery("pageIndex", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "5"))
	allRoom, err := global.GameSrvClient.SearchAllRoom(context.Background(), &game_proto.GetPageInfo{
		PageIndex: uint32(index),
		PageSize:  uint32(size),
	})
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
			"roomWait":      room.RoomWait,
			"roomName":      room.RoomName,
			"users":         user,
		})
		//fmt.Println(user)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": resp,
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

// 创建房间，做请求转发
func CreateRoom(ctx *gin.Context) {
	token := ctx.GetHeader("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "传入token!",
		})
		return
	}
	form := forms.CreateRoomForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	if len(global.ConsulProcessWebServices) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": "找不到目标服务器",
		})
		return
	}
	num := rand.Intn(len(global.ConsulProcessWebServices))
	//请求转发
	client := resty.New()
	resp, err := client.R().SetHeader("token", token).SetFormData(map[string]string{
		"room_id":         strconv.Itoa(form.RoomID),
		"max_user_number": strconv.Itoa(form.MaxUserNumber),
		"game_count":      strconv.Itoa(form.GameCount),
		"room_name":       form.RoomName,
	}).Post(fmt.Sprintf("http://%s:%d/v1/createRoom", global.ConsulProcessWebServices[num].ServiceAddress, global.ConsulProcessWebServices[num].ServicePort))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}
	ctx.JSON(resp.StatusCode(), string(resp.Body()))

	//ctx.Redirect(http.StatusPermanentRedirect,
	//	fmt.Sprintf("http://%s:%d/v1/createRoom", services[num].ServiceAddress, services[num].ServicePort))
}
