package initialize

import (
	"net/http"
	"process_web/middlewares"
	"process_web/router"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	engine := gin.Default()
	engine.GET("/health", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})
	//中间件
	engine.Use(middlewares.Cors())
	versionGroup := engine.Group("/v1")
	//使用中间件,要求只有登录的用户才能使用游戏接口
	router.InitCommonRouter(versionGroup)
	router.InitRoomRouter(versionGroup)
	return engine
}
