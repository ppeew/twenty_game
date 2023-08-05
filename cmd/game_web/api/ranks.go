package api

import (
	"context"
	"game_web/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

func GetRanks(ctx *gin.Context) {
	ranks, err := global.GameSrvClient.GetRanks(context.Background(), &emptypb.Empty{})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": ranks.Info,
	})
}
