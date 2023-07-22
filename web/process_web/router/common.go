package router

import (
	"process_web/api"
	"process_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitCommonRouter(group *gin.RouterGroup) {
	group.GET("connectSocket", middlewares.JWTAuthInParam(), api.ConnSocket)
}
