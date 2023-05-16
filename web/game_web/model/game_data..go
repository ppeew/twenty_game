package model

type Game struct {
	RoomID     uint32
	Users      map[uint32]*UserGameInfo
	GameCount  uint32
	UserNumber uint32

	InitChan   chan struct{}    //通知游戏初始化已经完成
	CommonChan chan Message     //游戏逻辑管道
	ChatChan   chan ChatMsgData //聊天管道
	ItemChan   chan ItemMsgData //使用物品管道
	HealthChan chan Message     //心脏包管道

	MakeCardID uint32 //依次生成卡的id
	RandCard   []Card //卡id->卡信息(包含特殊和普通卡)
}

type ChatMsgData struct {
	UserID uint32
	Data   []byte //聊天信息
}

type ItemMsgData struct {
	UserID       uint32
	Item         uint32 //使用的物品
	TargetUserID uint32
}

type Card struct {
	Type            uint32      `json:"type"`
	CardID          uint32      `json:"cardID"`
	SpecialCardInfo SpecialCard `json:"specialCardInfo"`
	BaseCardInfo    BaseCard    `json:"baseCardCardInfo"`
	HasOwner        bool        `json:"hasOwner"`
}

const (
	BaseType = iota
	SpecialType
)

type SpecialCard struct {
	CardID uint32
	Type   uint32
}

const (
	AddCard    = 1 << iota //增加卡
	DeleteCard             //删除卡
	UpdateCard             //更新卡
	ChangeCard             //交换卡
)

type BaseCard struct {
	CardID uint32
	Number uint32
}

type UserGameInfo struct {
	BaseCards    []BaseCard    //普通卡
	SpecialCards []SpecialCard //特殊卡
	IsGetCard    bool          //当前回合已经抢过卡了
	Score        uint32
}
