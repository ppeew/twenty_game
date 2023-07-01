package response

// 房间信息体
type RoomResponse struct {
	RoomID        uint32         `json:"roomID"`
	MaxUserNumber uint32         `json:"maxUserNumber"`
	GameCount     uint32         `json:"gameCount"`
	UserNumber    uint32         `json:"userNumber"`
	RoomOwner     uint32         `json:"roomOwner"`
	RoomWait      bool           `json:"roomWait"`
	RoomName      string         `json:"roomName"`
	Users         []UserResponse `json:"users"`
}

type UserResponse struct {
	ID    uint32 `json:"ID"`
	Ready bool   `json:"Ready"`
}

// 踢人的信息体，告知被t的用户
type KickerResponse struct {
}
