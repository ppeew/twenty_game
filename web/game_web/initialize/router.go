package initialize

import (
	"game_web/middlewares"
	"game_web/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	engine := gin.Default()
	engine.GET("/health", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})
	//中间件
	engine.Use(middlewares.FlowBegin(), middlewares.Cors())
	versionGroup := engine.Group("/v1")
	//使用中间件,要求只有登录的用户才能使用游戏接口
	router.InitCommonRouter(versionGroup)
	router.InitRoomRouter(versionGroup)
	router.InitGameRouter(versionGroup)
	router.InitRanksRouter(versionGroup)
	return engine
}
