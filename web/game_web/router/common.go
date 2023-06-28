package router

import (
	"game_web/api"
	"game_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitCommonRouter(group *gin.RouterGroup) {
	group.GET("reconnect", middlewares.JWTAuthInParam(), api.Reconnect, middlewares.FlowEnd())
	group.GET("getReconnInfo", middlewares.JWTAuth(), api.GetReconnInfo, middlewares.FlowEnd())
}
