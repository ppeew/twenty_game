package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"user_srv/global"
	"user_srv/model"
	"user_srv/proto/user"

	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	user.UnimplementedUserServer
}

// 用户注册
func (s *UserServer) CreateUser(ctx context.Context, req *user.CreateUserInfo) (*user.UserInfoResponse, error) {
	//先查询用户是否存在
	var u model.User
	result := global.MysqlDB.Where("user_name = ?", req.UserName).First(&u)
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

	u2 := &model.User{
		Nickname: req.Nickname,
		Gender:   req.Gender,
		UserName: req.UserName,
		Password: encodePassword,
		Good:     100,
		Diamond:  10,
	}
	res := global.MysqlDB.Create(u2)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return ModelToResponse(u2), nil
}

// 通过id获得用户信息
func (s *UserServer) GetUserByID(ctx context.Context, req *user.UserIDInfo) (*user.UserInfoResponse, error) {
	var u model.User
	result := global.MysqlDB.First(&u, req.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	return ModelToResponse(&u), nil
}

// 通过username获得用户信息
func (s *UserServer) GetUserByUsername(ctx context.Context, req *user.UserNameInfo) (*user.UserInfoResponse, error) {
	var u model.User
	result := global.MysqlDB.Where("user_name = ?", req.UserName).First(&u)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	return ModelToResponse(&u), nil
}

// 检查密码
func (s *UserServer) CheckPassword(ctx context.Context, req *user.CheckPasswordInfo) (*user.CheckPasswordResponse, error) {
	options := &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	info := strings.Split(req.EncodePassword, "$")
	verify := password.Verify(req.Password, info[0], info[1], options)
	return &user.CheckPasswordResponse{Success: verify}, nil
}

// 更改用户信息
func (s *UserServer) UpdateUser(ctx context.Context, req *user.UpdateUserInfo) (*emptypb.Empty, error) {
	var u model.User
	result := global.MysqlDB.First(&u, req.Id)
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
	if req.Password == "" {
		encodePassword = ""
	}
	res := global.MysqlDB.Model(&u).Where("id=?", fmt.Sprintf("%d", req.Id)).Updates(model.User{
		UserName: req.UserName,
		Password: encodePassword,
		Nickname: req.Nickname,
		Gender:   req.Gender,
	})
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.Internal, "更新用户失败")
	}
	return &emptypb.Empty{}, nil
}

func ModelToResponse(u *model.User) *user.UserInfoResponse {
	userInfoRep := &user.UserInfoResponse{
		Nickname: u.Nickname,
		Gender:   u.Gender,
		UserName: u.UserName,
		Password: u.Password,
		Id:       u.ID,
		Image:    u.Image,
	}
	return userInfoRep
}
