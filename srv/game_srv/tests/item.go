package main

import (
	"context"
	"fmt"
	"game_srv/proto"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
)

var userClient proto.GameClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("consul://192.168.159.134:8500/twelve_game_game_srv?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		panic(err)
	}
	userClient = proto.NewGameClient(conn)
}

func TestGetUserItem() {
	rsp, err := userClient.GetUserItemsInfo(context.Background(), &proto.UserIDInfo{
		Id: 1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id, rsp.Gold, rsp.Diamond)

}

func TestCreateUserItem() {
	rsp, err := userClient.CreateUserItems(context.Background(), &proto.UserItemsInfo{
		Id:      1,
		Gold:    10000,
		Diamond: 100,
		Apple:   0,
		Banana:  0,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func main() {
	Init()
	TestCreateUserItem()
	TestGetUserItem()

	conn.Close()
}
