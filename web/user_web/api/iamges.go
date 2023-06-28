package api

import (
	"context"
	"net/http"
	"strconv"
	"time"
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
	filePathByte, _ := time.Now().MarshalText()
	filePath := string(filePathByte) + "." + formFile.Filename
	err = ctx.SaveUploadedFile(formFile, filePath)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	_, err = global.UserSrvClient.UploadImage(context.Background(), &user.UploadInfo{Id: id, Path: filePath})
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
	ctx.JSON(http.StatusOK, image.Path)
}
