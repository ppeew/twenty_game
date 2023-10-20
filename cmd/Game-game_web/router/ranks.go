package router

import (
	"game_web/api"
	"game_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRanksRouter(group *gin.RouterGroup) {
	group.GET("getRanks", middlewares.Jaeger(), api.GetRanks, middlewares.FlowEnd())
}
