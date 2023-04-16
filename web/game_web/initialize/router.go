package initialize

import (
	"game_web/router"
	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	engine := gin.Default()
	//中间件
	//engine.Use()
	versionGroup := engine.Group("/v1")
	router.InitGameRouter(versionGroup)

	return engine
}
