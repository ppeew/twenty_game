package router

import (
	"game_web/api"

	"github.com/gin-gonic/gin"
)

func InitCommonRouter(group *gin.RouterGroup) {
	group.GET("selectUserState", api.SelectUserState)
}
