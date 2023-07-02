package router

import (
	"file_web/api"
	"file_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitFileRouter(r *gin.RouterGroup) {
	group := r.Group("file")
	group.Use(middlewares.JWTAuth())
	group.POST("uploadImage", api.UploadImage)
	group.GET("downloadImage", api.DownloadImage)
}
