package initialize

import (
	"game_web/middlewares"
	"game_web/router"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var num = int64(0)
var count = int64(5000)

func InitRouters() *gin.Engine {
	engine := gin.Default()
	engine.GET("/health", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})
	//中间件
	engine.Use(middlewares.FlowBegin(), middlewares.Cors())
	versionGroup := engine.Group("/v1")
	//使用中间件,要求只有登录的用户才能使用游戏接口
	isTest := true
	versionGroup.GET("/isStart", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"isTest": isTest,
			"begin":  num,
			"count":  count,
		})
		atomic.AddInt64(&num, count)
	})

	router.InitCommonRouter(versionGroup)
	router.InitRoomRouter(versionGroup)
	router.InitRanksRouter(versionGroup)
	return engine
}
