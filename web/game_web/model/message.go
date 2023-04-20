package model

// 用户与websocket通讯的消息体
type Message struct {
	UserID     uint32     `json:"userID"`
	Type       uint32     `json:"type"`
	DeleteData DeleteData `json:"deleteData"`
	UpdateData UpdateData `json:"updateData"`
	RoomData   RoomData   `json:"roomData"`
	ReadyState ReadyState `json:"readyState"`
	BeginGame  BeginGame  `json:"beginGame"`
}

const (
	DeleteRoom = 1 << iota
	UpdateRoom
	GetRoomData
	UserReadyState
	RoomBeginGame
)

type DeleteData struct {
	RoomID uint32 `json:"roomID"`
}
type UpdateData struct {
	MaxUserNumber uint32 `json:"maxUserNumber"`
	GameCount     uint32 `json:"gameCount"`
	Owner         uint32 `json:"owner"`
	Kicker        uint32 `json:"kicker"`
}
type RoomData struct {
	RoomID uint32 `json:"roomID"`
}
type ReadyState struct {
	IsReady bool `json:"isReady"`
}
type BeginGame struct {
	RoomID uint32 `json:"roomID"`
}
