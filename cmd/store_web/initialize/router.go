package initialize

import (
	"net/http"
	"store_web/middlewares"
	"store_web/service/router"

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
