package main

import (
	"store_srv/model"

	"github.com/spf13/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func main() {
	debug := GetEnvInfo("PPEEW_DEBUG")
	dsn := "root:518888@tcp(139.159.234.134:3306)/game?charset=utf8mb4&parseTime=True&loc=Local"
	if debug {
		dsn = "root:123456@tcp(127.0.0.1:3306)/twelve_game_store_srv?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&model.UserItem{})
}
