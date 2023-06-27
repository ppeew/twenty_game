package handler

import (
	"bufio"
	"context"
	"crypto/sha512"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"user_srv/global"
	"user_srv/model"
	"user_srv/proto/game"
	"user_srv/proto/user"

	"google.golang.org/grpc"

	"go.uber.org/zap"

	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	user.UnimplementedUserServer
}

// 上传头像文件
func (s *UserServer) UploadImage(ctx context.Context, in *user.UploadInfo) (*emptypb.Empty, error) {
	//先查询用户是否有了头像
	u := model.User{}
	first := global.MysqlDB.First(&u, in.Id)
	if first.RowsAffected != 1 {
		zap.S().Warnf("[UploadImage]:%s", first.Error)
		return &emptypb.Empty{}, first.Error
	}
	var filePath string
	if u.Image == "/" {
		//路径存储到数据库(没头像路径情况下)
		filePathByte, _ := time.Now().MarshalText()
		tx := global.MysqlDB.Model(&model.User{}).Where("id=?", in.Id).Update("image", filePathByte)
		if tx.Error != nil || tx.RowsAffected == 0 {
			zap.S().Warnf("[UploadImage]:%s", tx.Error)
			return &emptypb.Empty{}, tx.Error
		}
		filePath = "/usr/game_images/" + string(filePathByte)
	} else {
		//找到了头像，不需要修改数据库，将磁盘文件更改
		filePath = "/usr/game_images/" + u.Image
	}

	//写入服务器磁盘
	_ = os.MkdirAll("/usr/game_images/", 0666)
	openFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
	defer openFile.Close()
	if err != nil {
		zap.S().Warnf("[UploadImage]:%s", err)
		return &emptypb.Empty{}, err
	}
	writer := bufio.NewWriter(openFile)
	_, err = writer.Write(in.File)
	_ = writer.Flush()
	return &emptypb.Empty{}, nil
}

// 下载头像文件
func (s *UserServer) DownLoadImage(ctx context.Context, in *user.DownloadInfo) (*user.DownloadResponse, error) {
	//先查询用户是否有了头像
	u := model.User{}
	first := global.MysqlDB.First(&u, in.Id)
	if first.RowsAffected != 1 {
		zap.S().Warnf("[UploadImage]:%s", first.Error)
		return &user.DownloadResponse{}, first.Error
	}
	_ = os.MkdirAll("/usr/game_images/", 0666)
	filePath := "/usr/game_images/" + u.Image
	if u.Image == "/" {
		//没头像，用默认的
		filePath = "/usr/game_images/default.jpg"
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		zap.S().Warnf("[DownLoadImage]:%s", err)
		return &user.DownloadResponse{}, err
	}
	reader := bufio.NewReader(file)
	var ret []byte
	for true {
		p := make([]byte, 4096)
		n, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				//读到结尾
				break
			} else {
				return &user.DownloadResponse{}, err
			}
		}
		ret = append(ret, p[:n]...)
	}
	return &user.DownloadResponse{File: ret}, nil
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
	}
	tx := global.MysqlDB.Begin()
	res := tx.Create(u2)
	if res.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	if res.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	zap.S().Info("插入用户成功，接下来插入物品")
	//创建用户表成功，接下来为游戏用户添加物品表(跨服务调用,失败采用事务回滚)
	consulInfo := global.ServerConfig.ConsulInfo
	gameConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GameSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【游戏物品服务失败】")
	}
	client := game.NewGameClient(gameConn)
	_, err = client.CreateUserItems(ctx, &game.UserItemsInfo{
		Id:      u2.ID,
		Gold:    10000,
		Diamond: 100,
		Apple:   2,
		Banana:  2,
	})
	if err != nil {
		zap.S().Info(err.Error())
		tx.Rollback()
		return nil, err
	}
	zap.S().Info("插入物品成功，commit")
	tx.Commit()
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
	if result.Error != nil {
		return nil, result.Error
	}
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
	}
	return userInfoRep
}
