package utils

import (
	"hall_web/global"
	"hall_web/service/domains"
)

func CreateTable() {
	global.MysqlDB.AutoMigrate(&domains.Comments{})
}
