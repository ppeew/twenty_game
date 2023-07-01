package router

import (
	"game_web/api"
	"game_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitCommonRouter(group *gin.RouterGroup) {
	group.GET("connectSocket", middlewares.JWTAuthInParam(), api.ConnSocket, middlewares.FlowEnd())
	group.GET("getConnInfo", middlewares.JWTAuth(), api.GetConnInfo, middlewares.FlowEnd())
}
