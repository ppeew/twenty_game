package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game_srv/global"
	"game_srv/model"
	"game_srv/proto/game"
	"game_srv/utils"

	"go.uber.org/zap"
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

func (s *GameServer) CreateUserItems(ctx context.Context, req *game.UserItemsInfo) (*game.UserItemsInfoResponse, error) {
	zap.S().Info("用户访问CreateUserItems")
	fmt.Println("用户访问CreateUserItems")
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

func (s *GameServer) SearchAllRoom(ctx context.Context, in *emptypb.Empty) (*game.AllRoomInfo, error) {
	ret := &game.AllRoomInfo{}
	keys := global.RedisDB.Keys(ctx, "*")
	if keys.Err() != nil {
		return nil, keys.Err()
	}
	for _, value := range keys.Val() {
		get := global.RedisDB.Get(ctx, value)
		result := get.Val()
		if result == "" {
			return nil, get.Err()
		}
		room := model.Room{}
		_ = json.Unmarshal([]byte(result), &room)
		var users []*game.RoomUser
		for _, user := range room.Users {
			users = append(users, &game.RoomUser{
				ID:    user.ID,
				Ready: user.Ready,
			})
		}
		r := &game.RoomInfo{
			RoomID:        room.RoomID,
			MaxUserNumber: room.MaxUserNumber,
			GameCount:     room.GameCount,
			UserNumber:    room.UserNumber,
			RoomOwner:     room.RoomOwner,
			RoomWait:      room.RoomWait,
			Users:         users,
		}
		ret.AllRoomInfo = append(ret.AllRoomInfo, r)
	}
	return ret, nil
}

func (s *GameServer) CreateRoom(ctx context.Context, in *game.RoomInfo) (*emptypb.Empty, error) {
	var users []*model.User
	room := model.Room{
		RoomID:        in.RoomID,
		MaxUserNumber: in.MaxUserNumber,
		GameCount:     in.GameCount,
		UserNumber:    in.UserNumber,
		RoomOwner:     in.RoomOwner,
		RoomWait:      in.RoomWait,
		Users:         users,
	}
	marshal, _ := json.Marshal(room)
	set := global.RedisDB.Set(ctx, fmt.Sprintf("%d", in.RoomID), marshal, 0)
	if set.Err() != nil {
		return &emptypb.Empty{}, set.Err()
	}
	return &emptypb.Empty{}, nil
}

func (s *GameServer) SearchRoom(ctx context.Context, in *game.RoomIDInfo) (*game.RoomInfo, error) {
	get := global.RedisDB.Get(ctx, fmt.Sprintf("%d", in.RoomID))
	result := get.Val()
	if get.Val() == "" {
		return nil, get.Err()
	}
	room := model.Room{}
	_ = json.Unmarshal([]byte(result), &room)
	ret := &game.RoomInfo{
		RoomID:        room.RoomID,
		MaxUserNumber: room.MaxUserNumber,
		GameCount:     room.GameCount,
		UserNumber:    room.UserNumber,
		RoomOwner:     room.RoomOwner,
		RoomWait:      room.RoomWait,
		//Users:         make([]*game.RoomUser, 0),  不用初始化也可以append
	}
	for _, user := range room.Users {
		ret.Users = append(ret.Users, &game.RoomUser{
			ID:    user.ID,
			Ready: user.Ready,
		})
	}
	return ret, nil
}

func (s *GameServer) DeleteRoom(ctx context.Context, in *game.RoomIDInfo) (*emptypb.Empty, error) {
	del := global.RedisDB.Del(ctx, fmt.Sprintf("%d", in.RoomID))
	if del.Err() != nil {
		return &emptypb.Empty{}, del.Err()
	}
	if del.Val() == 0 {
		return &emptypb.Empty{}, errors.New("没有该房间")
	}
	return &emptypb.Empty{}, nil
}

func (s *GameServer) UserIntoRoom(ctx context.Context, in *game.UserIntoRoomInfo) (*game.IntoRoomRsp, error) {
	mutex, _ := utils.GetRedSync(in.RoomID)
	//查找房间是否存在
	get := global.RedisDB.Get(ctx, fmt.Sprintf("%d", in.RoomID))
	result := get.Val()
	if get.Val() == "" {
		_ = utils.ReleaseRedSync(mutex)
		return nil, get.Err()
	}
	room := model.Room{}
	_ = json.Unmarshal([]byte(result), &room)
	//是否满人
	if room.MaxUserNumber == room.UserNumber {
		_ = utils.ReleaseRedSync(mutex)
		return &game.IntoRoomRsp{ErrorMsg: "房间满人了"}, nil
	}
	//处理加入房间逻辑
	room.Users = append(room.Users, &model.User{
		ID:    in.UserID,
		Ready: false,
	})
	room.UserNumber++
	marshal, _ := json.Marshal(room)
	set := global.RedisDB.Set(ctx, fmt.Sprintf("%d", in.RoomID), marshal, 0)
	if set.Err() != nil {
		_ = utils.ReleaseRedSync(mutex)
		return nil, set.Err()
	}
	_ = utils.ReleaseRedSync(mutex)
	ret := &game.RoomInfo{
		RoomID:        room.RoomID,
		MaxUserNumber: room.MaxUserNumber,
		GameCount:     room.GameCount,
		UserNumber:    room.UserNumber,
		RoomOwner:     room.RoomOwner,
		RoomWait:      room.RoomWait,
	}
	for _, user := range room.Users {
		ret.Users = append(ret.Users, &game.RoomUser{
			ID:    user.ID,
			Ready: user.Ready,
		})
	}
	return &game.IntoRoomRsp{RoomInfo: ret}, nil
}

func (s *GameServer) QuitRoom(ctx context.Context, in *game.QuitRoomInfo) (*game.QuitRsp, error) {
	mutex, _ := utils.GetRedSync(in.RoomID)
	//查找房间是否存在
	get := global.RedisDB.Get(ctx, fmt.Sprintf("%d", in.RoomID))
	result := get.Val()
	if get.Val() == "" {
		_ = utils.ReleaseRedSync(mutex)
		return nil, get.Err()
	}
	room := model.Room{}
	_ = json.Unmarshal([]byte(result), &room)
	if in.UserID == room.RoomOwner {
		// 如果是房主退出，销毁房间
		del := global.RedisDB.Del(ctx, fmt.Sprintf("%d", in.RoomID))
		_ = utils.ReleaseRedSync(mutex)
		if del.Err() != nil {
			return nil, del.Err()
		}
		if del.Val() == 0 {
			return nil, errors.New("没有该房间")
		}
		ret := &game.RoomInfo{
			RoomID:        room.RoomID,
			MaxUserNumber: room.MaxUserNumber,
			GameCount:     room.GameCount,
			UserNumber:    room.UserNumber,
			RoomOwner:     room.RoomOwner,
			RoomWait:      room.RoomWait,
		}
		for _, user := range room.Users {
			ret.Users = append(ret.Users, &game.RoomUser{
				ID:    user.ID,
				Ready: user.Ready,
			})
		}
		return &game.QuitRsp{IsOwnerQuit: true, RoomInfo: ret}, nil
	} else {
		// 更新房间
		for i, user := range room.Users {
			if user.ID == in.UserID {
				room.Users = append(room.Users[:i], room.Users[i+1:]...)
				room.UserNumber--
			}
		}
		marshal, _ := json.Marshal(room)
		set := global.RedisDB.Set(ctx, fmt.Sprintf("%d", in.RoomID), marshal, 0)
		_ = utils.ReleaseRedSync(mutex)
		if set.Err() != nil {
			return nil, set.Err()
		}
		ret := &game.RoomInfo{
			RoomID:        room.RoomID,
			MaxUserNumber: room.MaxUserNumber,
			GameCount:     room.GameCount,
			UserNumber:    room.UserNumber,
			RoomOwner:     room.RoomOwner,
			RoomWait:      room.RoomWait,
		}
		for _, user := range room.Users {
			ret.Users = append(ret.Users, &game.RoomUser{
				ID:    user.ID,
				Ready: user.Ready,
			})
		}
		return &game.QuitRsp{IsOwnerQuit: false, RoomInfo: ret}, nil
	}
}

func (s *GameServer) UpdateRoom(ctx context.Context, in *game.UpdateRoomInfo) (*game.RoomInfo, error) {
	mutex, _ := utils.GetRedSync(in.RoomID)
	//查找房间是否存在
	get := global.RedisDB.Get(ctx, fmt.Sprintf("%d", in.RoomID))
	result := get.Val()
	if get.Val() == "" {
		_ = utils.ReleaseRedSync(mutex)
		return nil, get.Err()
	}
	room := model.Room{}
	_ = json.Unmarshal([]byte(result), &room)

	// 判断是不是房主
	if in.UserID != room.RoomOwner {
		_ = utils.ReleaseRedSync(mutex)
		return nil, errors.New("非房主不可修改")
	}

	if in.MaxUserNumber != 0 {
		room.MaxUserNumber = in.MaxUserNumber
	}
	if in.GameCount != 0 {
		room.GameCount = in.GameCount
	}
	if in.Owner != 0 {
		room.RoomOwner = in.Owner
	}
	if in.Kicker != 0 {
		if in.Kicker == room.RoomOwner {
			//不能t自己
			return nil, errors.New("不可T自己")
		}
		for i, user := range room.Users {
			if user.ID == in.Kicker {
				room.Users = append(room.Users[:i], room.Users[i+1:]...)
				room.UserNumber--
			}
		}
	}
	marshal, _ := json.Marshal(room)
	set := global.RedisDB.Set(ctx, fmt.Sprintf("%d", in.RoomID), marshal, 0)
	_ = utils.ReleaseRedSync(mutex)
	if set.Err() != nil {
		return nil, set.Err()
	}
	ret := &game.RoomInfo{
		RoomID:        room.RoomID,
		MaxUserNumber: room.MaxUserNumber,
		GameCount:     room.GameCount,
		UserNumber:    room.UserNumber,
		RoomOwner:     room.RoomOwner,
		RoomWait:      room.RoomWait,
	}
	for _, user := range room.Users {
		ret.Users = append(ret.Users, &game.RoomUser{
			ID:    user.ID,
			Ready: user.Ready,
		})
	}
	return ret, nil
}

func (s *GameServer) UpdateUserReadyState(ctx context.Context, in *game.ReadyStateInfo) (*game.RoomInfo, error) {
	mutex, _ := utils.GetRedSync(in.RoomID)
	//查找房间是否存在
	get := global.RedisDB.Get(ctx, fmt.Sprintf("%d", in.RoomID))
	result := get.Val()
	if get.Val() == "" {
		_ = utils.ReleaseRedSync(mutex)
		return nil, get.Err()
	}
	room := model.Room{}
	_ = json.Unmarshal([]byte(result), &room)
	for _, user := range room.Users {
		if user.ID == in.UserID {
			user.Ready = in.IsReady
		}
	}
	marshal, _ := json.Marshal(room)
	set := global.RedisDB.Set(ctx, fmt.Sprintf("%d", in.RoomID), marshal, 0)
	_ = utils.ReleaseRedSync(mutex)
	if set.Err() != nil {
		return nil, set.Err()
	}
	ret := &game.RoomInfo{
		RoomID:        room.RoomID,
		MaxUserNumber: room.MaxUserNumber,
		GameCount:     room.GameCount,
		UserNumber:    room.UserNumber,
		RoomOwner:     room.RoomOwner,
		RoomWait:      room.RoomWait,
	}
	for _, user := range room.Users {
		ret.Users = append(ret.Users, &game.RoomUser{
			ID:    user.ID,
			Ready: user.Ready,
		})
	}
	return ret, nil
}

func (s *GameServer) BeginGame(ctx context.Context, in *game.BeginGameInfo) (*game.BeginGameRsp, error) {
	mutex, _ := utils.GetRedSync(in.RoomID)
	//查找房间是否存在
	get := global.RedisDB.Get(ctx, fmt.Sprintf("%d", in.RoomID))
	result := get.Val()
	if get.Val() == "" {
		_ = utils.ReleaseRedSync(mutex)
		return nil, get.Err()
	}
	room := model.Room{}
	_ = json.Unmarshal([]byte(result), &room)
	//判断是否是房主
	if in.UserID != room.RoomOwner {
		_ = utils.ReleaseRedSync(mutex)
		return &game.BeginGameRsp{ErrorMsg: "非房主"}, nil
	}
	//检查是否满人
	if room.UserNumber != room.MaxUserNumber {
		_ = utils.ReleaseRedSync(mutex)
		return &game.BeginGameRsp{ErrorMsg: "没满人"}, nil
	}
	//检查其他人是否准备了
	ownerIndex := 0
	for i, user := range room.Users {
		if user.ID != room.RoomOwner && user.Ready == false {
			_ = utils.ReleaseRedSync(mutex)
			return &game.BeginGameRsp{ErrorMsg: "有玩家未准备"}, nil
		}
		if user.ID == room.RoomOwner {
			ownerIndex = i
		}
	}
	//可以开始游戏
	room.Users[ownerIndex].Ready = true
	room.RoomWait = false
	marshal, _ := json.Marshal(room)
	set := global.RedisDB.Set(ctx, fmt.Sprintf("%d", in.RoomID), marshal, 0)
	_ = utils.ReleaseRedSync(mutex)
	if set.Err() != nil {
		return nil, set.Err()
	}
	ret := &game.RoomInfo{
		RoomID:        room.RoomID,
		MaxUserNumber: room.MaxUserNumber,
		GameCount:     room.GameCount,
		UserNumber:    room.UserNumber,
		RoomOwner:     room.RoomOwner,
		RoomWait:      room.RoomWait,
	}
	for _, user := range room.Users {
		ret.Users = append(ret.Users, &game.RoomUser{
			ID:    user.ID,
			Ready: user.Ready,
		})
	}
	return &game.BeginGameRsp{RoomInfo: ret}, nil
}

func (s *GameServer) BackRoom(ctx context.Context, in *game.RoomIDInfo) (*emptypb.Empty, error) {
	mutex, _ := utils.GetRedSync(in.RoomID)
	//查找房间是否存在
	get := global.RedisDB.Get(ctx, fmt.Sprintf("%d", in.RoomID))
	result := get.Val()
	if get.Val() == "" {
		_ = utils.ReleaseRedSync(mutex)
		return nil, get.Err()
	}
	room := model.Room{}
	_ = json.Unmarshal([]byte(result), &room)

	// 更新房间状态为等待，所有玩家为未准备
	room.RoomWait = true
	for _, user := range room.Users {
		user.Ready = false
	}
	marshal, _ := json.Marshal(room)
	set := global.RedisDB.Set(ctx, fmt.Sprintf("%d", in.RoomID), marshal, 0)
	_ = utils.ReleaseRedSync(mutex)
	if set.Err() != nil {
		return nil, set.Err()
	}
	return &emptypb.Empty{}, nil
}
