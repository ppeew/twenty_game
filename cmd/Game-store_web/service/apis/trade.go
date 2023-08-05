package apis

import (
	"net/http"
	"store_web/model"
	"store_web/service/dao"
	"store_web/service/dto"

	"github.com/gin-gonic/gin"
)

func SelectTradeItems(ctx *gin.Context) {
	req := dto.TradeSelectReq{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	tradeDao := dao.TradeDao{}
	rsp, err := tradeDao.SelectPage(req.Page, req.Size)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}

func PushTradeItem(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	req := dto.TradePushReq{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	tradeDao := dao.TradeDao{}
	data, err := tradeDao.Create(req.ToDomain(int(userID)))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func DownTradeItem(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	req := dto.TradeDownReq{}
	err := ctx.BindUri(&req)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	tradeDao := dao.TradeDao{}
	err = tradeDao.DeleteByID(req.TradeItemID, userID)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusOK)
}

func BuyTradeItem(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	req := dto.TradeBuyReq{}
	err := ctx.BindUri(&req)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	tradeDao := dao.TradeDao{}

	err = tradeDao.BuyTradeItem(req.TradeItemID, int(userID))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusOK)
}
