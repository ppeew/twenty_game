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

// 获得用户重连需要的服务器信息（ip+port）
func (s *GameServer) GetConnData(ctx context.Context, in *game.UserIDInfo) (*game.ConnResponse, error) {
	//查询redis
	get := global.RedisDB.Get(ctx, NameUserConnInfo(in.Id))
	if get.Err() != nil {
		//找不到或者其他错误
		return &game.ConnResponse{}, errors.New("找不到对应的服务器")
	}
	return &game.ConnResponse{ServerInfo: get.Val()}, nil
}

// 记录用户连接的服务器信息
func (s *GameServer) RecordConnData(ctx context.Context, in *game.RecordConnInfo) (*emptypb.Empty, error) {
	mutex := global.RedSync.NewMutex(fmt.Sprintf("RecordConnLock:%d", in.Id))
	mutex.Lock()
	defer mutex.Unlock()
	get := global.RedisDB.Get(ctx, NameUserConnInfo(in.Id))
	if get.Err() != redis.Nil {
		//zap.S().Infof("[RecordConnData]:该用户已经有房间，不允许进房")
		return &emptypb.Empty{}, errors.New("该用户已经有房间，不允许进房")
	}
	//zap.S().Infof("[RecordConnData]:该用户在大厅，允许进房")
	set := global.RedisDB.Set(ctx, NameUserConnInfo(in.Id), in.ServerInfo, 0)
	if set.Err() != nil {
		return &emptypb.Empty{}, set.Err()
	}
	return &emptypb.Empty{}, nil
}

// 删除用户连接的服务器信息
func (s *GameServer) DelConnData(ctx context.Context, in *game.DelConnInfo) (*emptypb.Empty, error) {
	mutex := global.RedSync.NewMutex(fmt.Sprintf("RecordConnLock:%d", in.Id))
	mutex.Lock()
	defer mutex.Unlock()
	del := global.RedisDB.Del(ctx, NameUserConnInfo(in.Id))
	if del.Err() != nil {
		return &emptypb.Empty{}, del.Err()
	}
	return &emptypb.Empty{}, nil
}
