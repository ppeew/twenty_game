package response

// 给http客户端返回的信息体
type RoomResponse struct {
	RoomID        uint32         `json:"roomID"`
	MaxUserNumber uint32         `json:"maxUserNumber"`
	GameCount     uint32         `json:"gameCount"`
	UserNumber    uint32         `json:"userNumber"`
	RoomOwner     uint32         `json:"roomOwner"`
	RoomWait      bool           `json:"roomWait"`
	Users         []UserResponse `json:"users"`
}

type UserResponse struct {
	ID    uint32 `json:"ID"`
	Ready bool   `json:"Ready"`
}
