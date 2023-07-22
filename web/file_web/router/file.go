package router

import (
	"file_web/api"
	"file_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitFileRouter(r *gin.RouterGroup) {
	r.Use(middlewares.JWTAuth())
	r.POST("uploadImage", api.UploadImage)
	r.GET("downloadImage", api.DownloadImage)
}
