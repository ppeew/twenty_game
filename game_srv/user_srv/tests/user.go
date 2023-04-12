package main

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"user_srv/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("consul://192.168.159.134:8500/twenty_game_user_srv?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func TestGetUser() {
	rsp, err := userClient.GetUserByID(context.Background(), &proto.UserIDInfo{
		Id: 1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id, rsp.Name, rsp.OpenID, rsp.Gender)
	check, err := userClient.CheckOpenID(context.Background(), &proto.CheckOpenIDInfo{
		OpenID:       "admin123",
		EncodeOpenID: rsp.OpenID,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("检查openid", check.Success)

}

func TestCreateUser() {
	rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Name:   "bobby",
		OpenID: "admin123",
		Gender: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func main() {
	Init()
	//TestCreateUser()
	TestGetUser()

	conn.Close()
}
