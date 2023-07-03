package handler

import (
	"context"
	"fmt"
	"game_srv/global"
	"game_srv/proto/game"

	"github.com/redis/go-redis/v9"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 查询房间对应处理服务器信息
func (s *GameServer) GetRoomServer(ctx context.Context, in *game.RoomIDInfo) (*game.RoomServerInfo, error) {
	get := global.RedisDB.Get(ctx, NameRoomServer(in.RoomID))
	if get.Err() != nil {
		return &game.RoomServerInfo{}, get.Err()
	}
	return &game.RoomServerInfo{ServerInfo: get.Val()}, nil
}

// 删除房间对应处理服务器信息
func (s *GameServer) DelRoomServer(ctx context.Context, in *game.RoomIDInfo) (*emptypb.Empty, error) {
	del := global.RedisDB.Del(ctx, NameRoom(in.RoomID))
	return &emptypb.Empty{}, del.Err()
}

// 创建房间对应处理服务器信息
func (s *GameServer) RecordRoomServer(ctx context.Context, in *game.RecordRoomServerInfo) (*emptypb.Empty, error) {
	//先查询是否有了
	get := global.RedisDB.Get(ctx, NameRoomServer(in.RoomID))
	if get.Err() != redis.Nil {
		return &emptypb.Empty{}, get.Err()
	}
	global.RedisDB.Set(ctx, NameRoomServer(in.RoomID), in.ServerInfo, 0)
	return &emptypb.Empty{}, nil
}

func NameRoomServer(roomID uint32) string {
	return fmt.Sprintf("Game:RoomServer:%d", roomID)
}