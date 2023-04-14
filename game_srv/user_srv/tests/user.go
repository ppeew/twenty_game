package main

import (
	"context"
	"fmt"
	"user_srv/proto"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
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
	fmt.Println(rsp.Id, rsp.UserName, rsp.Password, rsp.Gender)

	rsp, err = userClient.GetUserByUsername(context.Background(), &proto.UserNameInfo{UserName: "root"})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id, rsp.UserName, rsp.Password, rsp.Gender)

	check, err := userClient.CheckPassword(context.Background(), &proto.CheckPasswordInfo{
		Password:       "123456",
		EncodePassword: rsp.Password,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("检查openid", check.Success)

}

func TestCreateUser() {
	rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Nickname: "ppeew",
		Gender:   true,
		UserName: "root",
		Password: "123456",
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
