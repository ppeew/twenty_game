package initialize

import (
	"file_web/middlewares"
	"file_web/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	engine := gin.Default()
	//中间件
	engine.Use(middlewares.Cors())
	engine.GET("/health", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})
	versionGroup := engine.Group("/v1")
	//router.InitOssRouter(versionGroup)
	router.InitFileRouter(versionGroup)
	return engine
}
