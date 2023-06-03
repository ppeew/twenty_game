package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestRedis(t *testing.T) {
	//Redis
	redisConn := redis.NewClient(&redis.Options{
		Addr:     "8.134.163.22:6379",
		DB:       0,
		Password: "ppeew",
	})
	get := redisConn.Get(context.Background(), "123")
	//找不到是这样的redis.Nil
	if get.Err() == redis.Nil {
		fmt.Println("find no")
		return
	}
	if get.Err() != nil {
		panic(get.Err())
	}
	fmt.Println(get.Val())
}
