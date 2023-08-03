package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game_srv/global"
	"game_srv/model"
	"game_srv/proto/game"

	"github.com/redis/go-redis/v9"

	"google.golang.org/protobuf/types/known/emptypb"
)

var RoomKey = "Game:Room"

// 查询所有房间
func (s *GameServer) SearchAllRoom(ctx context.Context, in *emptypb.Empty) (*game.AllRoomInfo, error) {
	ret := &game.AllRoomInfo{}
	//keys := global.RedisDB.Keys(ctx, "Game:roomID*")
	keys := global.RedisDB.HGetAll(ctx, RoomKey)
	if keys.Err() != nil {
		return nil, keys.Err()
	}
	for _, value := range keys.Val() {
		room := model.Room{}
		_ = json.Unmarshal([]byte(value), &room)
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
	//get := global.RedisDB.Get(ctx, NameRoom(in.RoomID))
	get := global.RedisDB.HGet(ctx, RoomKey, NameRoom(in.RoomID))
	if get.Err() == redis.Nil {
		return nil, errors.New("记录没找到")
	}
	result := get.Val()
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
	//global.RedisDB.Set(ctx, NameRoom(in.RoomID), marshal, 0)
	global.RedisDB.HSet(ctx, RoomKey, NameRoom(in.RoomID), marshal)
	return &emptypb.Empty{}, nil
}

// 删除房间
func (s *GameServer) DeleteRoom(ctx context.Context, in *game.RoomIDInfo) (*emptypb.Empty, error) {
	//global.RedisDB.Del(ctx, NameRoom(in.RoomID))
	global.RedisDB.HDel(ctx, RoomKey, NameRoom(in.RoomID))
	return &emptypb.Empty{}, nil
}

func NameRoom(roomID uint32) string {
	return fmt.Sprintf("roomID:%d", roomID)
}

func NameUserConnInfo(id uint32) string {
	return fmt.Sprintf("User:reconnInfo:%d", id)
}
