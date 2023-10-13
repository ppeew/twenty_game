package initialize

import (
	"context"
	"fmt"
	"game_web/global"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func InitDB() {
	// MongoDB
	mongoInfo := global.ServerConfig.MongoInfo
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d", mongoInfo.User, mongoInfo.Password, mongoInfo.Host, mongoInfo.Port)))
	if err != nil {
		zap.S().Fatalf("[InitDB]连接mongodb服务器错误:%s", err)
	}
	global.MongoDB = client.Database(mongoInfo.Database)
}
