package model

import (
	"context"
	"encoding/json"
	"errors"
	"game_web/global"
	"strconv"
)

type Room struct {
	RoomID        int               `json:"roomID"`
	MaxUserNumber int               `json:"maxUserNumber"`
	GameCount     int               `json:"gameCount"`
	UserNumber    int               `json:"userNumber"`
	RoomOwner     int               `json:"roomOwner"`
	RoomWait      bool              `json:"roomWait"`
	Users         map[int]*UserItem `json:"users"`

	//因为房间需要能够广播，设计订阅发布功能
	Publish *Publisher
}

func (r *Room) Create(room Room) (err error) {
	marshal, _ := json.Marshal(room)
	set := global.RedisDB.Set(context.Background(), strconv.Itoa(room.RoomID), marshal, 0)
	if set.Err() != nil {
		err = set.Err()
		return
	}
	return
}

func (r *Room) Select(room Room) (retRoom Room, err error) {
	get := global.RedisDB.Get(context.Background(), strconv.Itoa(room.RoomID))
	result := get.Val()
	if get.Val() == "nil" {
		err = get.Err()
		return
	}
	json.Unmarshal([]byte(result), &retRoom)
	return
}

func (r *Room) Delete(room Room) (err error) {
	del := global.RedisDB.Del(context.Background(), strconv.Itoa(room.RoomID))
	if del.Err() != nil {
		err = del.Err()
		return
	}
	if del.Val() == 0 {
		err = errors.New("没有该房间")
		return
	}
	return
}

func (r *Room) Update(room Room) (err error) {
	marshal, _ := json.Marshal(room)
	set := global.RedisDB.Set(context.Background(), strconv.Itoa(room.RoomID), marshal, 0)
	if set.Err() != nil {
		err = set.Err()
		return
	}
	return
}

// 处理房间信息
func (r *Room) HandleRoom() {
	//创建房间之后会调用该函数，创建独立协程，不断监听房间信息，并向订阅者发送相关的信息

}
