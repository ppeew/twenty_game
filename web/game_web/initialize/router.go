package initialize

import (
	"game_web/middlewares"
	"game_web/router"
	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	engine := gin.Default()
	//中间件
	//engine.Use()
	versionGroup := engine.Group("/v1")
	//使用中间件,要求只有登录的用户才能使用游戏接口
	versionGroup.Use(middlewares.JWTAuth())

	router.InitRoomRouter(versionGroup)
	return engine
}
