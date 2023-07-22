package handler

import (
	"context"
	"user_srv/global"
	"user_srv/model"
	"user_srv/proto/user"

	"go.uber.org/zap"
)

// 上传头像文件
func (s *UserServer) UploadImage(ctx context.Context, in *user.UploadInfo) (*user.UploadResponse, error) {
	//先查询用户是否有了头像
	u := model.User{}
	first := global.MysqlDB.First(&u, in.Id)
	if first.RowsAffected != 1 {
		zap.S().Warnf("[UploadImage]:%s", first.Error)
		return &user.UploadResponse{}, first.Error
	}
	//路径存储到数据库(没头像路径情况下)
	tx := global.MysqlDB.Model(&model.User{}).Where("id=?", in.Id).Update("image", in.Path)
	if tx.Error != nil || tx.RowsAffected == 0 {
		zap.S().Warnf("[UploadImage]:%s", tx.Error)
		return &user.UploadResponse{}, tx.Error
	}
	return &user.UploadResponse{Path: in.Path}, nil
}

// 下载头像文件
func (s *UserServer) DownLoadImage(ctx context.Context, in *user.DownloadInfo) (*user.DownloadResponse, error) {
	//先查询用户是否有了头像
	u := model.User{}
	first := global.MysqlDB.First(&u, in.Id)
	if first.RowsAffected != 1 {
		zap.S().Warnf("[UploadImage]:%s", first.Error)
		return &user.DownloadResponse{}, first.Error
	}
	filePath := u.Image
	return &user.DownloadResponse{Path: filePath}, nil
}
