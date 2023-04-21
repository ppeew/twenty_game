package model

type Game struct {
	RoomID     uint32
	Users      map[uint32]*UserGameInfo
	GameCount  uint32
	UserNumber uint32

	CommonChan chan Message //游戏逻辑管道
	ChatChan   chan ChatMsg //聊天管道
	ItemChan   chan ItemMsg //使用物品管道

	MakeCardID uint32 //依次生成卡的id
	RandCard   []Card //卡id->卡信息(包含特殊和普通卡)
}

type ChatMsg struct {
	UserID uint32
	Data   []byte //聊天信息
}

type ItemMsg struct {
	UserID       uint32
	Item         uint32 //使用的物品
	TargetUserID uint32
}

type Card struct {
	Type             uint32
	CardID           uint32
	SpecialCardInfo  SpecialCard
	BaseCardCardInfo BaseCard
	HasOwner         bool
}

const (
	SpecialType = 1 << iota
	BaseType
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
	BaseCards    []BaseCard    //普通卡（面向全部用户）
	SpecialCards []SpecialCard //特殊卡（仅自己可以看）
	Items        []uint32
	IsGetCard    bool //当前回合已经抢过卡了
	Score        uint32
	WS           *WSConn
}
