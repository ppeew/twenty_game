package router

import (
	"user_web/api"
	"user_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(group *gin.RouterGroup) {
	group.POST("register", api.UserRegister, middlewares.FlowEnd())
	group.POST("login", api.UserLogin, middlewares.FlowEnd())
	group.PUT("modify", middlewares.JWTAuth(), api.UserUpdate, middlewares.FlowEnd())
	group.GET("search", middlewares.JWTAuth(), api.GetUserInfo, middlewares.FlowEnd())
	group.GET("getNickname", api.GetRandomNickName, middlewares.FlowEnd())
	//group.GET("getUsername", api.GetRandomUsername, middlewares.FlowEnd())
	group.POST("uploadImage", middlewares.JWTAuth(), api.UploadImage)
	group.GET("downloadImage", middlewares.JWTAuth(), api.DownloadImage)
}
