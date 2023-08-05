package router

import (
	"process_web/api"
	"process_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoomRouter(group *gin.RouterGroup) {
	//房间相关接口
	group.POST("createRoom", middlewares.JWTAuth(), api.CreateRoom)
	group.PUT("userIntoRoom", middlewares.JWTAuth(), api.UserIntoRoom)
}
