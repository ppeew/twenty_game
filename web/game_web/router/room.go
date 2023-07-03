package router

import (
	"game_web/api"
	"game_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoomRouter(group *gin.RouterGroup) {
	//房间相关接口
	group.GET("getRoomList", middlewares.JWTAuth(), api.GetRoomList, middlewares.FlowEnd())
	group.PUT("userIntoRoom", middlewares.JWTAuth(), api.UserIntoRoom, middlewares.FlowEnd())
	group.GET("selectRoomServer", middlewares.JWTAuth(), api.SelectRoomServer, middlewares.FlowEnd())
}
