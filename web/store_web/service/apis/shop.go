package apis

import (
	"net/http"
	"store_web/model"
	"store_web/service/dao"
	"store_web/service/dto"

	"github.com/gin-gonic/gin"
)

// 分页查询商品
func SelectShopGoods(ctx *gin.Context) {
	req := dto.ShopSelectReq{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	shopDao := dao.ShopDao{}
	rsp, err := shopDao.SelectGoods(req)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}

// 购买商品
func BuyShopGood(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	req := dto.ShopBuyReq{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	shopDao := dao.ShopDao{}
	err = shopDao.BuyGood(req, int(userID))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusOK)
}
