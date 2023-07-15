package utils

import (
	"store_web/global"
	"store_web/service/domain"
)

func CreateTable() {
	global.MysqlDB.AutoMigrate(&domain.Good{})
	global.MysqlDB.AutoMigrate(&domain.TradeItem{})
	global.MysqlDB.AutoMigrate(&domain.Item{})
	global.MysqlDB.AutoMigrate(&domain.UserAndItem{})
	global.MysqlDB.AutoMigrate(&domain.User{})
}
