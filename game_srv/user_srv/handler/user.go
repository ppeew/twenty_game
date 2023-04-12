package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
	"user_srv/global"
	"user_srv/model"
	"user_srv/proto"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func ModelToResponse(user model.User) *proto.UserInfoResponse {
	userInfoRep := &proto.UserInfoResponse{
		Name:   user.Name,
		OpenID: user.OpenID,
		Gender: user.Gender,
	}
	return userInfoRep
}

func MakeEncodeOpenID(openID string) string {
	option := &password.Options{
		//对openid做加密
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	salt, encodePassword := password.Encode(openID, option)
	return fmt.Sprintf("%s$%s", salt, encodePassword)
}

// 用户注册
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	//先查询用户是否存在
	encodeOpenID := MakeEncodeOpenID(req.OpenID)
	var user model.User
	result := global.DB.Where("open_id = ?", encodeOpenID).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Error(codes.AlreadyExists, "用户已经存在")
	}
	user = model.User{
		Name:   req.Name,
		OpenID: encodeOpenID,
		Gender: req.Gender,
	}

	res := global.DB.Create(&user)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return ModelToResponse(user), nil
}

// 通过openid获得用户信息
func (s *UserServer) GetUserByOpenID(ctx context.Context, req *proto.UserOpenIDInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	encodeOpenID := MakeEncodeOpenID(req.OpenID)
	result := global.DB.Where("open_id = ?", encodeOpenID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	return ModelToResponse(user), nil
}

// 通过id获得用户信息
func (s *UserServer) GetUserByID(ctx context.Context, req *proto.UserIDInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	return ModelToResponse(user), nil
}

// openid验证密码
func (s *UserServer) CheckOpenID(ctx context.Context, req *proto.CheckOpenIDInfo) (*proto.CheckResponse, error) {
	options := &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	info := strings.Split(req.EncodeOpenID, "$")
	verify := password.Verify(req.OpenID, info[0], info[1], options)
	return &proto.CheckResponse{Success: verify}, nil
}

// 更改用户信息
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	res := global.DB.Model(&user).Where("id = ?", fmt.Sprintf("%d", req.Id)).Updates(map[string]interface{}{
		"name":   req.Name,
		"gender": req.Gender,
	})
	if res != nil {
		return nil, status.Error(codes.Internal, "更新用户失败")
	}
	return &emptypb.Empty{}, nil
}
