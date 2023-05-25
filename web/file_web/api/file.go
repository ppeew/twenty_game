package api

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "文件上传错误",
		})
		return
	}
	filePath := "./files/" + file.Filename
	err = ctx.SaveUploadedFile(file, filePath)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "文件上传错误",
		})
		return
	}
	oss, err := NewOSS()
	if err != nil {
		zap.S().Warnf("[UploadFile]:%s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": "文件上传错误",
		})
		return
	}
	objectKey := uuid.New().String()
	oss.UploadFile(objectKey, filePath)
	ctx.JSON(http.StatusOK, gin.H{
		"data": fmt.Sprintf("%s文件上传成功", file.Filename),
	})
	//用户头像库修改objectKey

}

func DownloadFile(ctx *gin.Context) {
	//从数据库找到该objectKey，在OSS服务查找
	//claims, _ := ctx.Get("claims")
	//userID := claims.(*model.CustomClaims).ID
	objectKey := ""

	oss, err := NewOSS()
	if err != nil {
		zap.S().Warnf("[UploadFile]:%s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": "文件下载错误",
		})
		return
	}
	oss.DownloadFile(objectKey, "")
	ctx.File("")
}
