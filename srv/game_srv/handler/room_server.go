package handler

import (
	"context"
	"errors"
	"fmt"
	"game_srv/global"
	"game_srv/proto/game"

	"github.com/redis/go-redis/v9"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 查询房间对应处理服务器信息
func (s *GameServer) GetRoomServer(ctx context.Context, in *game.RoomIDInfo) (*game.RoomServerInfo, error) {
	get := global.RedisDB.Get(ctx, NameRoomServer(in.RoomID))
	if get.Err() != nil || get.Err() == redis.Nil {
		return &game.RoomServerInfo{}, errors.New("找不到该房间")
	}
	return &game.RoomServerInfo{ServerInfo: get.Val()}, nil
}

// 删除房间对应处理服务器信息
func (s *GameServer) DelRoomServer(ctx context.Context, in *game.RoomIDInfo) (*emptypb.Empty, error) {
	mutex := global.RedSync.NewMutex(fmt.Sprintf("RoomServerLock:%d", in.RoomID))
	mutex.Lock()
	defer mutex.Unlock()
	del := global.RedisDB.Del(ctx, NameRoomServer(in.RoomID))
	return &emptypb.Empty{}, del.Err()
}

// 创建房间对应处理服务器信息
func (s *GameServer) RecordRoomServer(ctx context.Context, in *game.RecordRoomServerInfo) (*emptypb.Empty, error) {
	mutex := global.RedSync.NewMutex(fmt.Sprintf("RoomServerLock:%d", in.RoomID))
	mutex.Lock()
	defer mutex.Unlock()
	//先查询是否有了
	get := global.RedisDB.Get(ctx, NameRoomServer(in.RoomID))
	if get.Err() != redis.Nil {
		return &emptypb.Empty{}, errors.New("房间已存在了")
	}
	global.RedisDB.Set(ctx, NameRoomServer(in.RoomID), in.ServerInfo, 0)
	return &emptypb.Empty{}, nil
}

func NameRoomServer(roomID uint32) string {
	return fmt.Sprintf("Game:RoomServer:%d", roomID)
}
