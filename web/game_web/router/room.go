package router

import (
	"game_web/api"
	"github.com/gin-gonic/gin"
)

func InitRoomRouter(group *gin.RouterGroup) {
	//房间相关接口
	group.GET("sayHello", api.SayHello)
	group.GET("getRoomList", api.GetRoomList)
	group.GET("createRoom", api.CreateRoom)
	group.GET("dropRoom", api.DropRoom)
	group.GET("updateRoom", api.UpdateRoom)
	group.GET("userIntoRoom", api.UserIntoRoom)
	group.GET("roomInfo", api.RoomInfo)
	group.GET("userReady", api.UpdateUserReadyState)
	group.GET("beginGame", api.BeginGame)
}
