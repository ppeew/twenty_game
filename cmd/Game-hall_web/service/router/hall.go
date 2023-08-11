package router

import (
	"hall_web/middlewares"
	"hall_web/service/apis"

	"github.com/gin-gonic/gin"
)

func InitStoreRouter(r *gin.RouterGroup) {
	// 留言榜
	comment := r.Group("/comment")
	comment.GET("", middlewares.JWTAuth(), apis.CommentList)
	comment.POST("/add", middlewares.JWTAuth(), apis.AddComment)
	comment.PUT("/update", middlewares.JWTAuth(), apis.UpdateComment)
	comment.DELETE("/del/:id", middlewares.JWTAuth(), apis.DelComment)
	comment.PUT("/like", middlewares.JWTAuth(), apis.LikeComment)
}
