package router

import (
	"github.com/gin-gonic/gin"
	"user_web/api"
	"user_web/middlewares"
)

func InitUserRouter(group *gin.RouterGroup) {
	group.GET("register", api.UserRegister)
	group.GET("login", api.UserLogin)
	group.GET("modify", middlewares.JWTAuth(), api.UserUpdate)
	group.GET("search", middlewares.JWTAuth(), api.GetUserInfo)
	// add
}
