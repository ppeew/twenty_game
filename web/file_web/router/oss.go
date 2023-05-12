package router

import (
	"file_web/api"
	"file_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOssRouter(r *gin.RouterGroup) {
	group := r.Group("oss")
	group.Use(middlewares.JWTAuth())
	group.POST("uploadFile", api.UploadFile)
	group.GET("downloadFile", api.DownloadFile)
}
