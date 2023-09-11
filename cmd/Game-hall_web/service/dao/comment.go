package dao

import (
	"gorm.io/gorm"
	"hall_web/global"
	"hall_web/service/domains"
	"hall_web/service/dto"
)

type CommentDao struct {
}

func (CommentDao) CommentList(req *dto.CommentListReq) (res []domains.Comments, err error) {
	err = global.MysqlDB.Model(&domains.Comments{}).Offset((req.Page - 1) * req.Size).Limit(req.Size).Order("time desc").Find(&res).Error
	return
}

func (CommentDao) AddComment(req *domains.Comments) (err error) {
	err = global.MysqlDB.Create(req).Error
	return
}

func (CommentDao) UpdateComment(req *domains.Comments) (err error) {
	err = global.MysqlDB.Model(&domains.Comments{}).Where("user_id = ? and id = ?", req.UserID, req.ID).Updates(req).Error
	return
}

func (CommentDao) DelComment(req *domains.Comments) (err error) {
	err = global.MysqlDB.Delete(req).Error
	return
}

func (CommentDao) LikeComment(req *domains.Comments) (err error) {
	err = global.MysqlDB.Model(&domains.Comments{}).Where(req).Update("like", gorm.Expr("like + ?", 1)).Error
	return
}
