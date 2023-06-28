package handler

import (
	"context"
	"game_srv/global"
	"game_srv/proto/game"
)

// 获得用户重连需要的服务器信息（ip+port）
func (s *GameServer) GetReconnInfo(ctx context.Context, in *game.UserIDInfo) (*game.ReconnResponse, error) {
	//查询redis
	get := global.RedisDB.Get(ctx, NameUserReconnInfo(in.Id))
	if get.Err() != nil {
		//找不到或者其他错误
		return &game.ReconnResponse{}, get.Err()
	}
	return &game.ReconnResponse{ServerInfo: get.Val()}, nil
}
