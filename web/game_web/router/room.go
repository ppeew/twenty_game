package router

import (
	"game_web/api"
	"github.com/gin-gonic/gin"
)

func InitRoomRouter(group *gin.RouterGroup) {
	//房间相关接口
	group.GET("getRoomList", api.GetRoomList)
	group.GET("createRoom", api.CreateRoom)
}
