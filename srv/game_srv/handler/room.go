package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game_srv/global"
	"game_srv/model"
	"game_srv/proto/game"
	"game_srv/proto/user"
	"game_srv/utils"

	"github.com/redis/go-redis/v9"

	"github.com/go-redsync/redsync/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 0->大厅 1->房间 2->游戏
const (
	OutSide = iota
	RoomIn
	GameIn
)

// 查询所有房间
func (s *GameServer) SearchAllRoom(ctx context.Context, in *emptypb.Empty) (*game.AllRoomInfo, error) {
	ret := &game.AllRoomInfo{}
	keys := global.RedisDB.Keys(ctx, "Game:roomID*")
	if keys.Err() != nil {
		return nil, keys.Err()
	}
	for _, value := range keys.Val() {
		get := global.RedisDB.Get(ctx, value)
		result := get.Val()
		if get.Err() == redis.Nil {
			continue
		}
		room := model.Room{}
		_ = json.Unmarshal([]byte(result), &room)
		var users []*game.RoomUser
		for _, u := range room.Users {
			users = append(users, &game.RoomUser{
				ID:    u.ID,
				Ready: u.Ready,
			})
		}
		r := &game.RoomInfo{
			RoomID:        room.RoomID,
			MaxUserNumber: room.MaxUserNumber,
			GameCount:     room.GameCount,
			UserNumber:    room.UserNumber,
			RoomOwner:     room.RoomOwner,
			RoomWait:      room.RoomWait,
			RoomName:      room.RoomName,
			Users:         users,
		}
		ret.AllRoomInfo = append(ret.AllRoomInfo, r)
	}
	return ret, nil
}

// 查询某一房间
func (s *GameServer) SearchRoom(ctx context.Context, in *game.RoomIDInfo) (*game.RoomInfo, error) {
	get := global.RedisDB.Get(ctx, NameRoom(in.RoomID))
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
		RoomName:      room.RoomName,
		//Users:         make([]*game.RoomUser, 0),  不用初始化也可以append
	}
	for _, u := range room.Users {
		ret.Users = append(ret.Users, &game.RoomUser{
			ID:    u.ID,
			Ready: u.Ready,
		})
	}
	return ret, nil
}

// 创建房间
func (s *GameServer) CreateRoom(ctx context.Context, in *game.RoomInfo) (*emptypb.Empty, error) {
	get := global.RedisDB.Get(ctx, NameRoom(in.RoomID))
	if get.Err() != nil {
		if get.Err() != redis.Nil {
			return &emptypb.Empty{}, get.Err()
		}
	}
	//查看用户状态
	state, err := global.UserSrvClient.GetUserState(context.Background(), &user.UserIDInfo{Id: in.RoomOwner})
	if err != nil {
		zap.S().Warnf("[CreateRoom]:%s", err)
		return &emptypb.Empty{}, err
	}
	if state.State != OutSide {
		return &emptypb.Empty{}, err
	}
	if in.RoomID == 0 {
		return &emptypb.Empty{}, errors.New("无0房间")
	}
	var users []*model.User
	room := model.Room{
		RoomID:        in.RoomID,
		MaxUserNumber: in.MaxUserNumber,
		GameCount:     in.GameCount,
		UserNumber:    in.UserNumber,
		RoomOwner:     in.RoomOwner,
		RoomWait:      in.RoomWait,
		Users:         users,
		RoomName:      in.RoomName,
	}
	marshal, _ := json.Marshal(room)
	set := global.RedisDB.Set(ctx, NameRoom(in.RoomID), marshal, 0)
	if set.Err() != nil {
		return &emptypb.Empty{}, set.Err()
	}
	return &emptypb.Empty{}, nil
}

// 用户进入房间
func (s *GameServer) UserIntoRoom(ctx context.Context, in *game.UserIntoRoomInfo) (*game.IntoRoomRsp, error) {
	mutex, _ := utils.GetRedSync(in.RoomID)
	defer func(mutex *redsync.Mutex) {
		err := utils.ReleaseRedSync(mutex)
		if err != nil {
			zap.S().Errorf("[UserIntoRoom]%s", err)
		}
	}(mutex)
	//查找房间是否存在
	get := global.RedisDB.Get(ctx, NameRoom(in.RoomID))
	if get.Err() == redis.Nil {
		return nil, get.Err()
	}
	result := get.Val()
	room := model.Room{}
	_ = json.Unmarshal([]byte(result), &room)
	//是否满人
	if room.MaxUserNumber == room.UserNumber {
		return &game.IntoRoomRsp{ErrorMsg: "房间满人了"}, nil
	}
	state, err := global.UserSrvClient.GetUserState(context.Background(), &user.UserIDInfo{Id: in.UserID})
	if err != nil {
		zap.S().Warnf("[UserIntoRoom]:%s", err)
		return &game.IntoRoomRsp{}, err
	}
	switch state.State {
	case RoomIn:
		return &game.IntoRoomRsp{ErrorMsg: "请先退出之前的房间,再进入房间"}, nil
	case GameIn:
		return &game.IntoRoomRsp{ErrorMsg: "正在游戏中，请不要进房"}, nil
	}
	//处理加入房间逻辑
	room.Users = append(room.Users, &model.User{
		ID:    in.UserID,
		Ready: false,
	})
	room.UserNumber++
	marshal, _ := json.Marshal(room)
	set := global.RedisDB.Set(ctx, NameRoom(in.RoomID), marshal, 0)
	if set.Err() != nil {
		return nil, set.Err()
	}
	_, _ = global.UserSrvClient.UpdateUserState(context.Background(), &user.UpdateUserStateInfo{Id: in.UserID, State: RoomIn})
	ret := &game.RoomInfo{
		RoomID:        room.RoomID,
		MaxUserNumber: room.MaxUserNumber,
		GameCount:     room.GameCount,
		UserNumber:    room.UserNumber,
		RoomOwner:     room.RoomOwner,
		RoomWait:      room.RoomWait,
	}
	for _, u := range room.Users {
		ret.Users = append(ret.Users, &game.RoomUser{
			ID:    u.ID,
			Ready: u.Ready,
		})
	}
	return &game.IntoRoomRsp{RoomInfo: ret}, nil
}

// 退出房间
func (s *GameServer) QuitRoom(ctx context.Context, in *game.QuitRoomInfo) (*game.QuitRsp, error) {
	mutex, _ := utils.GetRedSync(in.RoomID)
	defer func(mutex *redsync.Mutex) {
		err := utils.ReleaseRedSync(mutex)
		if err != nil {
			zap.S().Errorf("[UserIntoRoom]%s", err)
		}
	}(mutex)
	//查找房间是否存在
	get := global.RedisDB.Get(ctx, NameRoom(in.RoomID))
	result := get.Val()
	if get.Val() == "" {
		return nil, get.Err()
	}
	room := model.Room{}
	_ = json.Unmarshal([]byte(result), &room)
	if in.UserID == room.RoomOwner {
		// 如果是房主退出，销毁房间
		del := global.RedisDB.Del(ctx, NameRoom(in.RoomID))
		if del.Err() != nil {
			return nil, del.Err()
		}
		if del.Val() == 0 {
			return nil, errors.New("没有该房间")
		}
		//更改全体用户状态
		for _, info := range room.Users {
			_, _ = global.UserSrvClient.UpdateUserState(context.Background(), &user.UpdateUserStateInfo{Id: info.ID, State: OutSide})
		}
		ret := &game.RoomInfo{
			RoomID:        room.RoomID,
			MaxUserNumber: room.MaxUserNumber,
			GameCount:     room.GameCount,
			UserNumber:    room.UserNumber,
			RoomOwner:     room.RoomOwner,
			RoomWait:      room.RoomWait,
		}
		for _, u := range room.Users {
			ret.Users = append(ret.Users, &game.RoomUser{
				ID:    u.ID,
				Ready: u.Ready,
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
		set := global.RedisDB.Set(ctx, NameRoom(in.RoomID), marshal, 0)
		if set.Err() != nil {
			return nil, set.Err()
		}
		_, _ = global.UserSrvClient.UpdateUserState(context.Background(), &user.UpdateUserStateInfo{Id: in.UserID, State: OutSide})
		ret := &game.RoomInfo{
			RoomID:        room.RoomID,
			MaxUserNumber: room.MaxUserNumber,
			GameCount:     room.GameCount,
			UserNumber:    room.UserNumber,
			RoomOwner:     room.RoomOwner,
			RoomWait:      room.RoomWait,
		}
		for _, u := range room.Users {
			ret.Users = append(ret.Users, &game.RoomUser{
				ID:    u.ID,
				Ready: u.Ready,
			})
		}
		return &game.QuitRsp{IsOwnerQuit: false, RoomInfo: ret}, nil
	}
}

// 修改房间
func (s *GameServer) UpdateRoom(ctx context.Context, in *game.UpdateRoomInfo) (*game.RoomInfo, error) {
	mutex, _ := utils.GetRedSync(in.RoomID)
	defer func(mutex *redsync.Mutex) {
		err := utils.ReleaseRedSync(mutex)
		if err != nil {
			zap.S().Errorf("[UserIntoRoom]%s", err)
		}
	}(mutex)
	//查找房间是否存在
	get := global.RedisDB.Get(ctx, NameRoom(in.RoomID))
	result := get.Val()
	if get.Val() == "" {
		return nil, get.Err()
	}
	room := model.Room{}
	_ = json.Unmarshal([]byte(result), &room)

	// 判断是不是房主
	if in.UserID != room.RoomOwner {
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
	set := global.RedisDB.Set(ctx, NameRoom(in.RoomID), marshal, 0)
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

// 更改用户状态
func (s *GameServer) UpdateUserReadyState(ctx context.Context, in *game.ReadyStateInfo) (*game.RoomInfo, error) {
	mutex, _ := utils.GetRedSync(in.RoomID)
	defer func(mutex *redsync.Mutex) {
		err := utils.ReleaseRedSync(mutex)
		if err != nil {
			zap.S().Errorf("[UserIntoRoom]%s", err)
		}
	}(mutex)
	//查找房间是否存在
	get := global.RedisDB.Get(ctx, NameRoom(in.RoomID))
	result := get.Val()
	if get.Val() == "" {
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
	set := global.RedisDB.Set(ctx, NameRoom(in.RoomID), marshal, 0)
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

// 开始游戏
func (s *GameServer) BeginGame(ctx context.Context, in *game.BeginGameInfo) (*game.BeginGameRsp, error) {
	mutex, _ := utils.GetRedSync(in.RoomID)
	defer func(mutex *redsync.Mutex) {
		err := utils.ReleaseRedSync(mutex)
		if err != nil {
			zap.S().Errorf("[UserIntoRoom]%s", err)
		}
	}(mutex)
	//查找房间是否存在
	get := global.RedisDB.Get(ctx, NameRoom(in.RoomID))
	result := get.Val()
	if get.Val() == "" {
		return nil, get.Err()
	}
	room := model.Room{}
	_ = json.Unmarshal([]byte(result), &room)
	//判断是否是房主
	if in.UserID != room.RoomOwner {
		return &game.BeginGameRsp{ErrorMsg: "非房主"}, nil
	}
	//检查是否满人
	if room.UserNumber != room.MaxUserNumber {
		return &game.BeginGameRsp{ErrorMsg: "没满人"}, nil
	}
	//检查其他人是否准备了
	ownerIndex := 0
	for i, u := range room.Users {
		if u.ID != room.RoomOwner && u.Ready == false {
			return &game.BeginGameRsp{ErrorMsg: "有玩家未准备"}, nil
		}
		if u.ID == room.RoomOwner {
			ownerIndex = i
		}
	}
	//可以开始游戏
	room.Users[ownerIndex].Ready = true
	room.RoomWait = false
	marshal, _ := json.Marshal(room)
	set := global.RedisDB.Set(ctx, NameRoom(in.RoomID), marshal, 0)
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
	for _, u := range room.Users {
		ret.Users = append(ret.Users, &game.RoomUser{
			ID:    u.ID,
			Ready: u.Ready,
		})
	}
	for _, info := range room.Users {
		_, _ = global.UserSrvClient.UpdateUserState(context.Background(), &user.UpdateUserStateInfo{Id: info.ID, State: GameIn})
	}
	return &game.BeginGameRsp{RoomInfo: ret}, nil
}

// 删除房间
func (s *GameServer) DeleteRoom(ctx context.Context, in *game.RoomIDInfo) (*emptypb.Empty, error) {
	get := global.RedisDB.Get(ctx, NameRoom(in.RoomID))
	result := get.Val()
	if get.Err() != redis.Nil {
		return &emptypb.Empty{}, get.Err()
	}
	room := model.Room{}
	_ = json.Unmarshal([]byte(result), &room)
	del := global.RedisDB.Del(ctx, fmt.Sprintf("%s", NameRoom(in.RoomID)))
	if del.Err() != nil {
		return &emptypb.Empty{}, del.Err()
	}
	if del.Val() == 0 {
		return &emptypb.Empty{}, errors.New("没有该房间")
	}
	for _, info := range room.Users {
		_, err := global.UserSrvClient.UpdateUserState(context.Background(), &user.UpdateUserStateInfo{Id: info.ID, State: OutSide})
		if err != nil {
			zap.S().Infof("[ReleaseResource]:%s", err)
		}
	}
	return &emptypb.Empty{}, nil
}

// 回到房间
func (s *GameServer) BackRoom(ctx context.Context, in *game.RoomIDInfo) (*emptypb.Empty, error) {
	mutex, _ := utils.GetRedSync(in.RoomID)
	defer func(mutex *redsync.Mutex) {
		err := utils.ReleaseRedSync(mutex)
		if err != nil {
			zap.S().Errorf("[UserIntoRoom]%s", err)
		}
	}(mutex)
	//查找房间是否存在
	get := global.RedisDB.Get(ctx, NameRoom(in.RoomID))
	result := get.Val()
	if get.Val() == "" {
		return nil, get.Err()
	}
	room := model.Room{}
	_ = json.Unmarshal([]byte(result), &room)
	// 更新房间状态为等待，所有玩家为未准备
	room.RoomWait = true
	for _, u := range room.Users {
		u.Ready = false
	}
	marshal, _ := json.Marshal(room)
	set := global.RedisDB.Set(ctx, NameRoom(in.RoomID), marshal, 0)
	if set.Err() != nil {
		return nil, set.Err()
	}
	for _, u := range room.Users {
		_, _ = global.UserSrvClient.UpdateUserState(context.Background(), &user.UpdateUserStateInfo{Id: u.ID, State: RoomIn})
	}
	return &emptypb.Empty{}, nil
}

func NameRoom(roomID uint32) string {
	return fmt.Sprintf("Game:roomID:%d", roomID)
}
