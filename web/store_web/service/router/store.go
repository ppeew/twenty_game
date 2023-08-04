package router

import (
	"store_web/middlewares"
	"store_web/service/apis"

	"github.com/gin-gonic/gin"
)

func InitStoreRouter(r *gin.RouterGroup) {
	trade := r.Group("/trade")
	trade.GET("", middlewares.JWTAuth(), apis.SelectTradeItems)
	trade.POST("", middlewares.JWTAuth(), apis.PushTradeItem)
	trade.PUT("/:id", middlewares.JWTAuth(), apis.BuyTradeItem)
	trade.DELETE("/:id", middlewares.JWTAuth(), apis.DownTradeItem)

	shop := r.Group("/shop")
	shop.GET("", middlewares.JWTAuth(), apis.SelectShopGoods)
	shop.PUT("/:id", middlewares.JWTAuth(), apis.BuyShopGood)
}
