package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"user_srv/model"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/twenty_game_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&model.User{})
}
