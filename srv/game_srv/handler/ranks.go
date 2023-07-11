package handler

import (
	"context"
	"fmt"
	"game_srv/global"
	"game_srv/proto/game"
	"strconv"

	"go.uber.org/zap"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/redis/go-redis/v9"
)

var ranks string = "Game:ranks"

// 获得排行榜信息
func (s *GameServer) GetRanks(ctx context.Context, in *emptypb.Empty) (*game.RanksResponse, error) {
	//1.用redis的zset直接取出排行榜内容，zrange key 获取到排行榜key
	zRange := global.RedisDB.ZRevRange(ctx, ranks, 0, 100)
	if zRange.Err() == redis.Nil {
		zap.S().Info("[GetRanks]:%s", zRange.Err())
		return &game.RanksResponse{}, nil
	}
	strings := zRange.Val()
	var info []*game.UserRankInfo
	for _, s := range strings {
		id, _ := strconv.Atoi(s)
		scoreCmd := global.RedisDB.Get(ctx, NameUserScore(uint32(id)))
		score, _ := strconv.Atoi(scoreCmd.Val())
		gametimesCmd := global.RedisDB.Get(ctx, NameUserGametimes(uint32(id)))
		gametimes, _ := strconv.Atoi(gametimesCmd.Val())
		info = append(info, &game.UserRankInfo{
			Id:        uint32(id),
			Score:     uint64(score),
			Gametimes: uint64(gametimes),
		})
	}
	//2.整理返回排行榜信息，应该包含排名+相关游戏信息（总得分+总游戏次数）
	return &game.RanksResponse{Info: info}, nil
}

// 更新排行榜
func (s *GameServer) UpdateRanks(ctx context.Context, in *game.UpdateRanksInfo) (*emptypb.Empty, error) {
	mutex := global.RedSync.NewMutex(fmt.Sprintf("User:Ranks:%d", in.UserID))
	mutex.Lock()
	defer mutex.Unlock()
	//1.从redis中变更用户的信息（包含总得分+总游戏次数）
	times := global.RedisDB.IncrBy(ctx, NameUserGametimes(in.UserID), int64(in.AddGametimes))
	score := global.RedisDB.IncrBy(ctx, NameUserScore(in.UserID), int64(in.AddScore))
	//2.给redis同步数据更新,且更新排行榜zset的key
	global.RedisDB.ZAdd(ctx, ranks, redis.Z{
		Score:  float64(score.Val()) / float64(times.Val()),
		Member: in.UserID, //key:数字，userID
	})
	return &emptypb.Empty{}, nil
}

func NameUserScore(userID uint32) string {
	return fmt.Sprintf("Game:userScore:%d", userID)
}

func NameUserGametimes(userID uint32) string {
	return fmt.Sprintf("Game:userGametimes:%d", userID)
}
