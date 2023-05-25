package middlewares

import (
	"net/http"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
)

func FlowBegin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		e, b := sentinel.Entry("user_web", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"err": "请求过于频繁，请稍后重试",
			})
			ctx.Abort()
			return
		}
		ctx.Set("flow", e)
	}
}

func FlowEnd() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		flow, _ := ctx.Get("flow")
		e := flow.(*base.SentinelEntry)
		e.Exit()
	}
}
