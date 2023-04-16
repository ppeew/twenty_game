package initialize

import (
	"fmt"
	"game_web/global"
	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	//初始化redis
	redisInfo := global.ServerConfig.RedisInfo
	global.RedisDB = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", redisInfo.Host, redisInfo.Port),
		DB:   redisInfo.Database,
	})
}
