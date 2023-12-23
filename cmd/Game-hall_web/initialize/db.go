package initialize

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hall_web/global"
	"hall_web/service/domains"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	//Mysql
	mysqlInfo := global.ServerConfig.MysqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlInfo.User, mysqlInfo.Password, mysqlInfo.Host, mysqlInfo.Port, mysqlInfo.Database)
	zap.S().Infof("[InitDB]:dsn=%s", dsn)
	var err error
	global.MysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.S().Fatalf("[InitDB]打开mysql错误:%s", err.Error())
	}
	if err := global.MysqlDB.AutoMigrate(&domains.Comments{}); err != nil {
		zap.S().Infof("[InitDB]:%s", err)
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

	// MongoDB
	mongoInfo := global.ServerConfig.MongoInfo
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d", mongoInfo.User, mongoInfo.Password, mongoInfo.Host, mongoInfo.Port)))
	if err != nil {
		zap.S().Fatalf("[InitDB]连接mongodb服务器错误:%s", err)
	}
	global.MongoDB = client.Database(mongoInfo.Database)
}
