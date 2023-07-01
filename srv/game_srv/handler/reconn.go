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
	get := global.RedisDB.Get(ctx, NameUserReconnInfo(in.Id))
	if get.Err() != nil {
		//找不到或者其他错误
		return &game.ConnResponse{}, get.Err()
	}
	return &game.ConnResponse{ServerInfo: get.Val()}, nil
}

// 记录连接的服务器信息
func (s *GameServer) RecordConnData(ctx context.Context, in *game.RecordConnInfo) (*emptypb.Empty, error) {
	set := global.RedisDB.Set(ctx, NameUserReconnInfo(in.Id), in.ServerInfo, 0)
	if set.Err() != nil {
		return &emptypb.Empty{}, set.Err()
	}
	return &emptypb.Empty{}, nil
}
