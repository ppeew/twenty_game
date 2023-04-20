package initialize

import (
	"github.com/gin-gonic/gin"
	"user_web/middlewares"
	"user_web/router"
)

func InitRouters() *gin.Engine {
	engine := gin.Default()
	//中间件
	engine.Use(middlewares.Cors())
	versionGroup := engine.Group("/v1")
	router.InitUserRouter(versionGroup)

	return engine
}
