package router

import (
	"game_web/api"
	"game_web/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRoomRouter(group *gin.RouterGroup) {
	//房间相关接口
	group.GET("getRoomList", middlewares.JWTAuth(), api.GetRoomList)
	group.GET("createRoom", middlewares.JWTAuth(), api.CreateRoom)
	group.GET("getRoomInfo", middlewares.JWTAuth(), api.GetRoomInfo)
	group.GET("userIntoRoom", middlewares.JWTAuthInParam(), api.UserIntoRoom)
}
