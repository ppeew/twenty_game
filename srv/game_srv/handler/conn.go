package handler

import (
	"context"
	"game_srv/global"
	"game_srv/proto/game"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 获得用户重连需要的服务器信息（ip+port）
func (s *GameServer) GetConnData(ctx context.Context, in *game.UserIDInfo) (*game.ConnResponse, error) {
	//查询redis
	get := global.RedisDB.Get(ctx, NameUserConnInfo(in.Id))
	if get.Err() != nil {
		//找不到或者其他错误
		return &game.ConnResponse{}, get.Err()
	}
	return &game.ConnResponse{ServerInfo: get.Val()}, nil
}

// 记录用户连接的服务器信息
func (s *GameServer) RecordConnData(ctx context.Context, in *game.RecordConnInfo) (*emptypb.Empty, error) {
	set := global.RedisDB.Set(ctx, NameUserConnInfo(in.Id), in.ServerInfo, 0)
	if set.Err() != nil {
		return &emptypb.Empty{}, set.Err()
	}
	return &emptypb.Empty{}, nil
}

// 删除用户连接的服务器信息
func (s *GameServer) DelConnData(ctx context.Context, in *game.DelConnInfo) (*emptypb.Empty, error) {
	del := global.RedisDB.Del(ctx, NameUserConnInfo(in.Id))
	if del.Err() != nil {
		return &emptypb.Empty{}, del.Err()
	}
	return &emptypb.Empty{}, nil
}
