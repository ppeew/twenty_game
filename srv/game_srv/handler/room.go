package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game_srv/global"
	"game_srv/model"
	"game_srv/proto/game"
	"time"

	"github.com/redis/go-redis/v9"

	"google.golang.org/protobuf/types/known/emptypb"
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
func (s *GameServer) SetGlobalRoom(ctx context.Context, in *game.RoomInfo) (*emptypb.Empty, error) {
	if in.RoomID <= 0 {
		return &emptypb.Empty{}, errors.New("房间号不能小于0")
	}
	get := global.RedisDB.Get(ctx, NameRoom(in.RoomID))
	if get.Err() != nil {
		if get.Err() != redis.Nil {
			return &emptypb.Empty{}, errors.New("房间已存在了")
		}
		//房间不存在，允许创房
	}

	var users []*model.User
	for _, user := range in.Users {
		users = append(users, &model.User{ID: user.ID, Ready: user.Ready})
	}
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
	global.RedisDB.Set(ctx, NameRoom(in.RoomID), marshal, time.Second*10)
	//新建连接的redis服务器信息
	//global.RedisDB.Set(ctx, NameUserConnInfo(in.RoomOwner), fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port), 0)
	return &emptypb.Empty{}, nil
}

// 删除房间
func (s *GameServer) DeleteRoom(ctx context.Context, in *game.RoomIDInfo) (*emptypb.Empty, error) {
	get := global.RedisDB.Get(ctx, NameRoom(in.RoomID))
	if get.Err() == redis.Nil {
		return &emptypb.Empty{}, get.Err()
	}
	room := model.Room{}
	_ = json.Unmarshal([]byte(get.Val()), &room)
	for _, info := range room.Users {
		global.RedisDB.Del(ctx, NameUserConnInfo(info.ID)) //删除了用户连接信息
	}
	global.RedisDB.Del(ctx, fmt.Sprintf("%s", NameRoom(in.RoomID)))
	return &emptypb.Empty{}, nil
}

func NameRoom(roomID uint32) string {
	return fmt.Sprintf("Game:roomID:%d", roomID)
}

func NameUserConnInfo(id uint32) string {
	return fmt.Sprintf("User:reconnInfo:%d", id)
}
