package initialize

import (
	"context"
	"fmt"
	"user_srv/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	//Mysql
	mysqlInfo := global.ServerConfig.MysqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlInfo.User, mysqlInfo.Password, mysqlInfo.Host, mysqlInfo.Port, mysqlInfo.Database)
	var err error
	global.MysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.S().Fatalf("[InitDB]打开mysql错误:%s", err.Error())
	}

	//Redis
	redisInfo := global.ServerConfig.RedisInfo
	global.RedisDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisInfo.Host, redisInfo.Port),
		DB:       redisInfo.Database,
		Password: redisInfo.Password,
	})
	err = global.RedisDB.Ping(context.Background()).Err()
	if err != nil {
		zap.S().Fatalf("[InitDB]连接redis服务器错误:%s", err)
	}
}
