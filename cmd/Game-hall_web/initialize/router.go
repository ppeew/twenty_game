package initialize

import (
	"hall_web/middlewares"
	"hall_web/service/router"
	"net/http"

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
	router.InitStoreRouter(versionGroup)
	return engine
}
