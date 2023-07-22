package router

import (
	"store_web/service/apis"

	"github.com/gin-gonic/gin"
)

func InitStoreRouter(r *gin.RouterGroup) {
	trade := r.Group("/trade")
	trade.GET("", apis.SelectTradeItems)
	trade.POST("/:id", apis.PushTradeItem)
	trade.PUT("/:id", apis.BuyTradeItem)
	trade.DELETE("/:id", apis.DownTradeItem)

	shop := r.Group("/shop")
	shop.GET("", apis.SelectShopGoods)
	shop.PUT("/:id", apis.BuyShopGood)
}
