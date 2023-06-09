package tests

import (
	"context"
	"fmt"
	"testing"
	"user_srv/proto/user"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
)

var userClient user.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("consul://192.168.159.134:8500/twelve_game_user_srv?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		panic(err)
	}
	userClient = user.NewUserClient(conn)
}

func TestGetUser(t *testing.T) {
	rsp, err := userClient.GetUserByID(context.Background(), &user.UserIDInfo{
		Id: 1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id, rsp.UserName, rsp.Password, rsp.Gender)

	rsp, err = userClient.GetUserByUsername(context.Background(), &user.UserNameInfo{UserName: "root"})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id, rsp.UserName, rsp.Password, rsp.Gender)

	check, err := userClient.CheckPassword(context.Background(), &user.CheckPasswordInfo{
		Password:       "123456",
		EncodePassword: rsp.Password,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("检查openid", check.Success)

}

func TestCreateUser(t *testing.T) {
	rsp, err := userClient.CreateUser(context.Background(), &user.CreateUserInfo{
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
