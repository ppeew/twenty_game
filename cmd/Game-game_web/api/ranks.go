package api

import (
	"context"
	"game_web/global"
	"game_web/proto/game"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BindPage struct {
	PageIndex uint32 `form:"pageIndex"`
	PageSize  uint32 `form:"pageSize"`
}

func GetRanks(ctx *gin.Context) {
	var pageReq BindPage
	ctx.BindQuery(&pageReq)
	ranks, err := global.GameSrvClient.GetRanks(context.Background(), &game.GetPageInfo{PageSize: pageReq.PageSize, PageIndex: pageReq.PageIndex})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": ranks.Info,
	})
}
