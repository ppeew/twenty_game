package handler

import (
	"context"
	"fmt"
	"game_srv/global"
	"game_srv/model"
	"game_srv/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GameServer struct {
	proto.UnimplementedGameServer
}

func ModelToResponse(user *model.UserItem) *proto.UserItemsInfoResponse {
	var record []uint32
	record = append(record, user.Apple)
	record = append(record, user.Banana)
	userInfoRep := &proto.UserItemsInfoResponse{
		Id:      user.ID,
		Gold:    user.Gold,
		Diamond: user.Diamond,
		Items:   record,
	}
	return userInfoRep
}

func (s *GameServer) CreateUserItems(ctx context.Context, req *proto.UserItemsInfo) (*proto.UserItemsInfoResponse, error) {
	zap.S().Info("用户访问CreateUserItems")
	item := model.UserItem{
		Gold:    req.Gold,
		Diamond: req.Diamond,
		Apple:   req.Apple,
		Banana:  req.Banana,
	}
	if req.Id != 0 {
		item.BaseModel.ID = req.Id
	}
	res := global.MysqlDB.Create(&item)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}
	return ModelToResponse(&item), nil
}

func (s *GameServer) GetUserItemsInfo(ctx context.Context, req *proto.UserIDInfo) (*proto.UserItemsInfoResponse, error) {
	item := model.UserItem{
		BaseModel: model.BaseModel{
			ID: req.Id,
		},
	}
	res := global.MysqlDB.First(&item)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	return ModelToResponse(&item), nil
}

// 增加金币
func (s *GameServer) AddGold(ctx context.Context, req *proto.AddGoldRequest) (*emptypb.Empty, error) {
	res := global.MysqlDB.Model(&model.UserItem{}).Where("id = ?", fmt.Sprintf("%d", req.Id)).Update("gold = ?", fmt.Sprintf("%d", req.Count))
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.Internal, "更新用户失败")
	}
	return &emptypb.Empty{}, nil
}

// 增加钻石
func (s *GameServer) AddDiamond(ctx context.Context, req *proto.AddDiamondInfo) (*emptypb.Empty, error) {
	res := global.MysqlDB.Model(&model.UserItem{}).Where("id = ?", fmt.Sprintf("%d", req.Id)).Update("diamond = ?", fmt.Sprintf("%d", req.Count))
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.Internal, "更新用户失败")
	}
	return &emptypb.Empty{}, nil
}

// 增加道具(道具类型应该区别)
func (s *GameServer) AddItem(ctx context.Context, req *proto.AddItemInfo) (*emptypb.Empty, error) {
	//要知道更新什么
	item := model.UserItem{
		Apple:  req.Items[proto.Type_Apple],
		Banana: req.Items[proto.Type_Banana],
	}
	res := global.MysqlDB.Model(&model.UserItem{}).Where("id = ?", fmt.Sprintf("%d", req.Id)).Updates(item)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.Internal, "更新用户失败")
	}
	return &emptypb.Empty{}, nil
}
