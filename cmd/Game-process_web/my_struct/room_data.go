package my_struct

import "time"

type UserRoomData struct {
	IntoRoomTime time.Time
	ID           uint32
	Ready        bool
	Nickname     string `json:"nickname"`
	Gender       bool   `json:"gender"`
	Username     string `json:"username"`
	Image        string `json:"image"`
}

// 退出房间结构体
type QuitRoomData struct {
}

// 更新房间结构体
type UpdateData struct {
	MaxUserNumber uint32 `json:"maxUserNumber"`
	GameCount     uint32 `json:"gameCount"`
	RoomName      string `json:"roomName"`
	Owner         uint32 `json:"owner"`
	Kicker        uint32 `json:"kicker"`
}

// 查询房间数据结构体
type RoomData struct {
}

// 更新用户准备状态结构体
type ReadyStateData struct {
	IsReady bool `json:"isReady"`
}

// 开始游戏结构体
type BeginGameData struct {
}
