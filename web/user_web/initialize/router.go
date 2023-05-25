package initialize

import (
	"user_web/middlewares"
	"user_web/router"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	engine := gin.Default()
	//中间件
	engine.Use(middlewares.FlowBegin(), middlewares.Cors())
	versionGroup := engine.Group("/v1")
	router.InitUserRouter(versionGroup)
	return engine
}
