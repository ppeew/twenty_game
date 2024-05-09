package utils

import (
	"user_srv/global"
	"user_srv/model"
)

func CreateTable() {
	global.MysqlDB.AutoMigrate(&model.User{})
}
