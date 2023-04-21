package model

// 用户与websocket通讯的消息体
type Message struct {
	//通用
	Type    uint32  `json:"type"`
	ChatMsg ChatMsg `json:"chatMsgInfo"`
	//房间
	UserID     uint32     `json:"userID"`
	DeleteData DeleteData `json:"deleteData"`
	UpdateData UpdateData `json:"updateData"`
	RoomData   RoomData   `json:"roomData"`
	ReadyState ReadyState `json:"readyState"`
	BeginGame  BeginGame  `json:"beginGame"`
	//游戏
	ItemMsg        ItemMsg        `json:"itemMsgInfo"`
	GetCardData    GetCardData    `json:"getCardData"`
	UseSpecialData UseSpecialData `json:"useSpecialData"`
}

const (
	// 通用的消息
	ChatMsgData = 1 << iota
	// 房间相关的消息
	DeleteRoom
	UpdateRoom
	GetRoomData
	UserReadyState
	RoomBeginGame
	// 游戏相关的消息
	ItemMsgData
	ListenHandleCard
	UseSpecialCard
)

type UseSpecialData struct {
	SpecialCardID uint32 `json:"specialCardID"`
	//增加卡需要指定生成的数字卡点数
	AddCardData AddCardData `json:"addCardData"`
	//更改卡需要指定目标一张数字卡牌，变成什么（1-11）
	UpdateCardData UpdateCardData `json:"updateCardData"`
	//删除卡需要删除指定玩家的一张数字卡
	DeleteCardData DeleteCardData `json:"deleteCardData"`
	//交换卡需要指定自己的一张数字卡和对方玩家ID的一张数字卡
	ChangeCardData ChangeCardData `json:"changeCardData"`
}

// 交换卡结构体
type ChangeCardData struct {
	CardID       uint32 `json:"cardID"`
	TargetUserID uint32 `json:"targetUserID"`
	TargetCard   uint32 `json:"targetCard"`
}

// 删除卡结构体
type DeleteCardData struct {
	TargetUserID uint32 `json:"targetUserID"`
	CardID       uint32 `json:"cardID"`
}

// 更改卡结构体
type UpdateCardData struct {
	TargetUserID uint32 `json:"targetUserID"`
	CardID       uint32 `json:"cardID"`
	UpdateNumber uint32 `json:"updateNumber"`
}

// 增加卡结构体
type AddCardData struct {
	NeedNumber uint32 `json:"needNumber"`
}

// 抢卡结构体
type GetCardData struct {
	GetCardID uint32 `json:"getCardID"`
}

// 删除房间结构体
type DeleteData struct {
	RoomID uint32 `json:"roomID"`
}

// 更新房间结构体
type UpdateData struct {
	MaxUserNumber uint32 `json:"maxUserNumber"`
	GameCount     uint32 `json:"gameCount"`
	Owner         uint32 `json:"owner"`
	Kicker        uint32 `json:"kicker"`
}

// 查询房间数据结构体
type RoomData struct {
	RoomID uint32 `json:"roomID"`
}

// 更新装备状态结构体
type ReadyState struct {
	IsReady bool `json:"isReady"`
}

// 开始游戏结构体
type BeginGame struct {
	RoomID uint32 `json:"roomID"`
}
