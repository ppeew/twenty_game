package router

import (
	"game_web/api"
	"game_web/middlewares"
	"github.com/gin-gonic/gin"
)

func InitGameRouter(group *gin.RouterGroup) {
	//使用中间件,要求只有登录的用户才能使用游戏接口
	group.Use(middlewares.JWTAuth())
	group.GET("sayHello", api.SayHello)
	group.GET("getRoomList", api.GetRoomList)
	group.GET("createRoom", api.CreateRoom)
	group.GET("dropRoom", api.DropRoom)
	group.GET("updateRoom", api.UpdateRoom)
}
