package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"user_srv/global"
	"user_srv/model"
	"user_srv/proto"

	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func ModelToResponse(user *model.User) *proto.UserInfoResponse {
	userInfoRep := &proto.UserInfoResponse{
		Nickname: user.Nickname,
		Gender:   user.Gender,
		UserName: user.UserName,
		Password: user.Password,
		Id:       user.ID,
	}
	return userInfoRep
}

// 用户注册
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	//先查询用户是否存在
	var user model.User
	result := global.DB.Where("user_name = ?", req.UserName).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Error(codes.AlreadyExists, "用户已经存在")
	}

	options := &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	salt, encoded := password.Encode(req.Password, options)
	encodePassword := fmt.Sprintf("%s$%s", salt, encoded)

	user = model.User{
		Nickname: req.Nickname,
		Gender:   req.Gender,
		UserName: req.UserName,
		Password: encodePassword,
	}

	res := global.DB.Create(&user)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return ModelToResponse(&user), nil
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
	return ModelToResponse(&user), nil
}

func (s *UserServer) GetUserByUsername(ctx context.Context, req *proto.UserNameInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where("user_name = ?", req.UserName).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	return ModelToResponse(&user), nil
}

func (s *UserServer) CheckPassword(ctx context.Context, req *proto.CheckPasswordInfo) (*proto.CheckPasswordResponse, error) {
	options := &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	info := strings.Split(req.EncodePassword, "$")
	verify := password.Verify(req.Password, info[0], info[1], options)
	return &proto.CheckPasswordResponse{Success: verify}, nil
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

	options := &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	salt, encoded := password.Encode(req.Password, options)
	encodePassword := fmt.Sprintf("%s$%s", salt, encoded)

	res := global.DB.Model(&user).Where("id = ?", fmt.Sprintf("%d", req.Id)).Updates(map[string]interface{}{
		"nickname":  req.Nickname,
		"gender":    req.Gender,
		"user_name": req.UserName,
		"password":  encodePassword,
	})
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.Internal, "更新用户失败")
	}
	return &emptypb.Empty{}, nil
}
