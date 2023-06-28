package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"user_web/global"
	"user_web/models"
	"user_web/proto/user"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UploadImage(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	id := currentUser.ID
	formFile, err := ctx.FormFile("image")
	if err != nil {
		zap.S().Infof("[UploadImage]:%s", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	file, err := formFile.Open()
	if err != nil {
		zap.S().Infof("[UploadImage]:%s", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	data, _ := ioutil.ReadAll(file)
	_, err = global.UserSrvClient.UploadImage(context.Background(), &user.UploadInfo{File: data, Id: id})
	if err != nil {
		zap.S().Infof("[UploadImage]:%s", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusOK)
}

func DownloadImage(ctx *gin.Context) {
	idStr := ctx.DefaultQuery("id", "0")
	id, _ := strconv.Atoi(idStr)
	if id == 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}
	image, err := global.UserSrvClient.DownLoadImage(context.Background(), &user.DownloadInfo{Id: uint32(id)})
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, image.File)
}
