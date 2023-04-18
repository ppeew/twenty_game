package response

// 给客户端返回的信息体
type RoomResponse struct {
	RoomID        int                  `json:"roomID"`
	MaxUserNumber int                  `json:"maxUserNumber"`
	GameCount     int                  `json:"gameCount"`
	UserNumber    int                  `json:"userNumber"`
	RoomOwner     int                  `json:"roomOwner"`
	RoomWait      bool                 `json:"roomWait"`
	Users         map[int]UserResponse `json:"users"`
}
