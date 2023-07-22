package tests

import (
	"bufio"
	"context"
	"os"
	"testing"
	"user_srv/proto/user"

	"google.golang.org/grpc"
)

func TestImage(t *testing.T) {
	conn, err := grpc.Dial("consul://8.134.163.22:8500/twelve_game_user_srv?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		panic(err)
	}
	userClient = user.NewUserClient(conn)
	_, err = userClient.UploadImage(context.Background(), &user.UploadInfo{File: []byte("666666"), Id: 10})
	if err != nil {
		panic(err)
	}
	println("OK~")
	file, err := userClient.DownLoadImage(context.Background(), &user.DownloadInfo{Id: 5})
	if err != nil {
		panic(err)
	}
	open, _ := os.OpenFile("C:\\Users\\22378\\GolandProjects\\twenty_game\\srv\\user_srv\\images\\test.jpg", os.O_CREATE, 0666)
	writer := bufio.NewWriter(open)
	writer.Write(file.File)
	open.Close()
}
