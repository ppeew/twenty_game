package router

import (
	"user_web/api"
	"user_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(group *gin.RouterGroup) {
	group.POST("register", api.UserRegister)
	group.POST("login", api.UserLogin)
	group.PUT("modify", middlewares.JWTAuth(), api.UserUpdate)
	group.GET("search", api.GetUserInfo)
	group.GET("getNickname", api.GetRandomNickName)
	//group.GET("getUsername", api.GetRandomUsername, middlewares.FlowEnd())
	//group.POST("uploadImage", middlewares.JWTAuth(), api.UploadImage)
	//group.GET("downloadImage", middlewares.JWTAuth(), api.DownloadImage)
}
