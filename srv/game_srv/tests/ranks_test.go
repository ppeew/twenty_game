package tests

import (
	"context"
	"fmt"
	game_proto "game_srv/proto/game"
	"testing"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc"
)

func TestRanks(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.Dial("consul://8.134.163.22:8500/twelve_game_game_srv?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		panic(err)
	}
	client := game_proto.NewGameClient(conn)
	ranks, err := client.GetRanks(ctx, &emptypb.Empty{})
	if err != nil {
		panic(err)
	}
	for i, info := range ranks.Info {
		fmt.Println("第", i, "名:", info.Id, info.Score, info.Gametimes)
	}
}

func TestUpdateRanks(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.Dial("consul://8.134.163.22:8500/twelve_game_game_srv?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		panic(err)
	}
	client := game_proto.NewGameClient(conn)
	_, err = client.UpdateRanks(ctx, &game_proto.UpdateRanksInfo{
		UserID:       3,
		AddScore:     100,
		AddGametimes: 1,
	})
	if err != nil {
		panic(err)
	}
}
