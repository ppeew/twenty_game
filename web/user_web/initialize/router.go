package initialize

import (
	"net/http"
	"user_web/middlewares"
	"user_web/router"

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
	router.InitUserRouter(versionGroup)
	return engine
}
