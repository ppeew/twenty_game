package response

import "time"

// 房间信息体
type RoomResponse struct {
	RoomID        uint32     `json:"roomID"`
	MaxUserNumber uint32     `json:"maxUserNumber"`
	GameCount     uint32     `json:"gameCount"`
	UserNumber    uint32     `json:"userNumber"`
	RoomOwner     uint32     `json:"roomOwner"`
	RoomWait      bool       `json:"roomWait"`
	RoomName      string     `json:"roomName"`
	Users         []UserData `json:"users"`
}

type UserData struct {
	ID           uint32    `json:"ID"`
	Ready        bool      `json:"Ready"`
	IntoRoomTime time.Time `json:"-"` //忽略
	Nickname     string    `json:"nickname"`
	Gender       bool      `json:"gender"`
	Username     string    `json:"username"`
	Image        string    `json:"image"`
}

// 踢人的信息体，告知所有用户是谁被t了
type KickerResponse struct {
	ID uint32 `json:"ID"`
}

type BeginGameData struct {
}
