package api

//
//func UploadImage(ctx *gin.Context) {
//	claims, _ := ctx.Get("claims")
//	currentUser := claims.(*models.CustomClaims)
//	id := currentUser.ID
//	formFile, err := ctx.FormFile("image")
//	if err != nil {
//		zap.S().Infof("[UploadImage]:%s", err)
//		ctx.Status(http.StatusBadRequest)
//		return
//	}
//	filePathByte, _ := time.Now().MarshalText()
//	filePath := string(filePathByte)
//	filePath = strings.Split(filePath, ".")[0]
//	filePath = strings.Replace(filePath, "T", "/", 1) + "-" + formFile.Filename
//	prePath := "/opt/images/"
//	err = ctx.SaveUploadedFile(formFile, prePath+filePath)
//	if err != nil {
//		ctx.Status(http.StatusInternalServerError)
//		return
//	}
//	_, err = global.UserSrvClient.UploadImage(context.Background(), &user.UploadInfo{Id: id, Path: filePath})
//	if err != nil {
//		zap.S().Infof("[UploadImage]:%s", err)
//		ctx.Status(http.StatusBadRequest)
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"data": "OK",
//	})
//}
//
//func DownloadImage(ctx *gin.Context) {
//	idStr := ctx.DefaultQuery("id", "0")
//	id, _ := strconv.Atoi(idStr)
//	if id == 0 {
//		ctx.Status(http.StatusBadRequest)
//		return
//	}
//	image, err := global.UserSrvClient.DownLoadImage(context.Background(), &user.DownloadInfo{Id: uint32(id)})
//	if err != nil {
//		ctx.Status(http.StatusBadRequest)
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"path": image.Path,
//	})
//}
