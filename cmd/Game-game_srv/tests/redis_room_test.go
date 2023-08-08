package tests

import (
	"context"
	"fmt"
	game_proto "game_srv/proto/game"
	"google.golang.org/grpc"
	"testing"
)

func TestUpdateRedisRoom(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.Dial("consul://139.159.234.134:8500/game-srv-dev?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		panic(err)
	}
	client := game_proto.NewGameClient(conn)
	client.SetGlobalRoom(ctx, &game_proto.RoomInfo{
		RoomID:        331,
		MaxUserNumber: 4,
		GameCount:     1,
		UserNumber:    1,
		RoomOwner:     666,
		RoomWait:      true,
		Users:         nil,
		RoomName:      "test_1",
	})
	if err != nil {
		panic(err)
	}
	room, err := client.SearchAllRoom(ctx, &game_proto.GetPageInfo{
		PageIndex: 2,
		PageSize:  2,
	})
	fmt.Println(err)
	fmt.Println(room)

	searchRoom, _ := client.SearchRoom(ctx, &game_proto.RoomIDInfo{RoomID: 333})
	fmt.Println(searchRoom)

	_, err = client.DeleteRoom(ctx, &game_proto.RoomIDInfo{RoomID: 33})
	fmt.Println(err)
}
