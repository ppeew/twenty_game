package handler

import (
	"context"
	"user_srv/proto"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

// 用户注册
func CreateUser(context.Context, *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	
}

// 用户登录
func UserLogin(context.Context, *proto.UserLoginInfo) (*proto.UserInfoResponse, error) {

}

// 通过id获得用户信息
func GetUserById(context.Context, *proto.UserIdInfo) (*proto.UserInfoResponse, error) {

}

// 更改信息
func UpdateUser(context.Context, *proto.UpdateUserInfo) (*proto.UserInfoResponse, error) {

}