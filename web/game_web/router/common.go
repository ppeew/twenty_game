package router

import (
	"game_web/api"
	"game_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitCommonRouter(group *gin.RouterGroup) {
	group.GET("reconnect", middlewares.JWTAuthInParam(), api.Reconnect, middlewares.FlowEnd())
	group.GET("selectUserState", middlewares.JWTAuth(), api.SelectUserState, middlewares.FlowEnd())
}
