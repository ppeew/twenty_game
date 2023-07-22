package handler

import (
	"context"
	"errors"
	"fmt"
	"game_srv/global"
	"game_srv/model"
	"game_srv/proto/game"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GameServer struct {
	game.UnimplementedGameServer
}

func ModelToResponse(user *model.UserItem) *game.UserItemsInfoResponse {
	var record []uint32
	record = append(record, user.Apple)
	record = append(record, user.Banana)
	userInfoRep := &game.UserItemsInfoResponse{
		Id:      user.ID,
		Gold:    user.Gold,
		Diamond: user.Diamond,
		Items:   record,
	}
	return userInfoRep
}

// 创建用户物品表
func (s *GameServer) CreateUserItems(ctx context.Context, req *game.UserItemsInfo) (*game.UserItemsInfoResponse, error) {
	//zap.S().Info("用户访问CreateUserItems")
	item := model.UserItem{
		UserID:  req.Id,
		Gold:    req.Gold,
		Diamond: req.Diamond,
		Apple:   req.Apple,
		Banana:  req.Banana,
	}
	res := global.MysqlDB.Create(&item)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}
	return ModelToResponse(&item), nil
}

// 获得用户物品表
func (s *GameServer) GetUserItemsInfo(ctx context.Context, req *game.UserIDInfo) (*game.UserItemsInfoResponse, error) {
	item := model.UserItem{
		UserID: req.Id,
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
func (s *GameServer) AddGold(ctx context.Context, req *game.AddGoldRequest) (*emptypb.Empty, error) {
	query := &model.UserItem{
		UserID: req.Id,
	}
	tx := global.MysqlDB.First(&query)
	if tx.Error != nil {
		return &emptypb.Empty{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	query.Gold += req.Count
	res := global.MysqlDB.Model(&model.UserItem{}).Where("id = ?", fmt.Sprintf("%d", req.Id)).Update("gold = ?", fmt.Sprintf("%d", query.Gold))
	if res.RowsAffected == 0 {
		return &emptypb.Empty{}, status.Error(codes.Internal, "更新用户失败")
	}
	return &emptypb.Empty{}, nil
}

// 增加钻石
func (s *GameServer) AddDiamond(ctx context.Context, req *game.AddDiamondInfo) (*emptypb.Empty, error) {
	query := &model.UserItem{
		UserID: req.Id,
	}
	tx := global.MysqlDB.First(&query)
	if tx.Error != nil {
		return &emptypb.Empty{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	query.Diamond += req.Count
	res := global.MysqlDB.Model(&model.UserItem{}).Where("id = ?", fmt.Sprintf("%d", req.Id)).Update("diamond = ?", fmt.Sprintf("%d", query.Diamond))
	if res.RowsAffected == 0 {
		return &emptypb.Empty{}, status.Error(codes.Internal, "更新用户失败")
	}
	return &emptypb.Empty{}, nil
}

// 增加道具(道具类型应该区别)
func (s *GameServer) AddItem(ctx context.Context, req *game.AddItemInfo) (*emptypb.Empty, error) {
	//要知道更新什么
	query := &model.UserItem{
		UserID: req.Id,
	}
	tx := global.MysqlDB.First(&query)
	if tx.Error != nil {
		return &emptypb.Empty{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	item := model.UserItem{
		Apple:  req.Items[game.Type_Apple] + query.Apple,
		Banana: req.Items[game.Type_Banana] + query.Banana,
	}
	res := global.MysqlDB.Model(&model.UserItem{}).Where("id = ?", fmt.Sprintf("%d", req.Id)).Updates(item)
	if res.RowsAffected == 0 {
		return &emptypb.Empty{}, status.Error(codes.Internal, "更新用户失败")
	}
	return &emptypb.Empty{}, nil
}

// 使用金币
func (s *GameServer) UseGold(ctx context.Context, req *game.UseGoldRequest) (*game.IsOK, error) {
	//要知道更新什么
	query := &model.UserItem{
		UserID: req.Id,
	}
	tx := global.MysqlDB.First(&query)
	if tx.Error != nil {
		return &game.IsOK{IsOK: false}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return &game.IsOK{IsOK: false}, status.Error(codes.NotFound, "用户不存在")
	}
	if req.Count > query.Gold {
		//不可以使用
		return &game.IsOK{IsOK: false}, errors.New("道具不足，无法使用")
	}
	query.Gold -= req.Count
	res := global.MysqlDB.Model(&model.UserItem{}).Where("id = ?", fmt.Sprintf("%d", req.Id)).Update("diamond = ?", fmt.Sprintf("%d", query.Gold))
	if res.RowsAffected == 0 {
		return &game.IsOK{IsOK: false}, status.Error(codes.Internal, "更新用户失败")
	}
	return &game.IsOK{IsOK: true}, nil

}

// 使用钻石
func (s *GameServer) UseDiamond(ctx context.Context, req *game.UseDiamondInfo) (*game.IsOK, error) {
	//要知道更新什么
	query := &model.UserItem{
		UserID: req.Id,
	}
	tx := global.MysqlDB.First(&query)
	if tx.Error != nil {
		return &game.IsOK{IsOK: false}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return &game.IsOK{IsOK: false}, status.Error(codes.NotFound, "用户不存在")
	}
	if req.Count > query.Diamond {
		//不可以使用
		return &game.IsOK{IsOK: false}, errors.New("道具不足，无法使用")
	}
	query.Diamond -= req.Count
	res := global.MysqlDB.Model(&model.UserItem{}).Where("id = ?", fmt.Sprintf("%d", req.Id)).Update("diamond = ?", fmt.Sprintf("%d", query.Diamond))
	if res.RowsAffected == 0 {
		return &game.IsOK{IsOK: false}, status.Error(codes.Internal, "更新用户失败")
	}
	return &game.IsOK{IsOK: true}, nil
}

// 使用道具
func (s *GameServer) UseItem(ctx context.Context, req *game.UseItemInfo) (*game.IsOK, error) {
	//要知道更新什么
	query := &model.UserItem{
		UserID: req.Id,
	}
	tx := global.MysqlDB.First(&query)
	if tx.Error != nil {
		return &game.IsOK{IsOK: false}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return &game.IsOK{IsOK: false}, status.Error(codes.NotFound, "用户不存在")
	}
	item := model.UserItem{
		Apple:  query.Apple - req.Items[game.Type_Apple],
		Banana: query.Banana - req.Items[game.Type_Banana],
	}
	if item.Apple < 0 || item.Banana < 0 {
		//不可以使用
		return &game.IsOK{IsOK: false}, errors.New("道具不足，无法使用")
	}
	res := global.MysqlDB.Model(&model.UserItem{}).Where("id = ?", fmt.Sprintf("%d", req.Id)).Updates(item)
	if res.RowsAffected == 0 {
		return &game.IsOK{IsOK: false}, status.Error(codes.Internal, "更新用户失败")
	}
	return &game.IsOK{IsOK: true}, nil
}
