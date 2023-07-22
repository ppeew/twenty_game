package utils

import (
	"store_web/global"
	"store_web/service/domains"
)

func CreateTable() {
	global.MysqlDB.AutoMigrate(&domains.Good{})
	global.MysqlDB.AutoMigrate(&domains.TradeItem{})
	global.MysqlDB.AutoMigrate(&domains.Item{})
	global.MysqlDB.AutoMigrate(&domains.UserAndItem{})
	global.MysqlDB.AutoMigrate(&domains.User{})
}
