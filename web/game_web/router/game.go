package router

import (
	"game_web/api"
	"game_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitGameRouter(group *gin.RouterGroup) {
	group.GET("selectItems", middlewares.JWTAuth(), api.SelectItems)
}
