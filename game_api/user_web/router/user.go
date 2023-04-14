package router

import (
	"user_web/api"
	"user_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(group *gin.RouterGroup) {
	group.POST("register", api.UserRegister)
	group.POST("login", api.UserLogin)
	group.POST("modify", middlewares.JWTAuth(), api.UserUpdate)
	group.GET("search", middlewares.JWTAuth(), api.GetUserInfo)
	// add
}
