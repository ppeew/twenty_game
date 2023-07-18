package my_struct

import "time"

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
	ID               uint32         `json:"ID,omitempty"`
	BaseCards        []*BaseCard    `json:"baseCards,omitempty"`    //普通卡
	SpecialCards     []*SpecialCard `json:"specialCards,omitempty"` //特殊卡
	GetBaseCardNum   int            `json:"isGetCard,omitempty"`    //当前回合抢普通卡数量
	IsGetSpecialCard bool           `json:"isGetSpecialCard"`       //当前回合已经抢过特殊卡了
	Score            uint32         `json:"score,omitempty"`
	IntoRoomTime     time.Time      `json:"-"`
	Nickname         string         `json:"nickname"`
	Gender           bool           `json:"gender"`
	Username         string         `json:"username"`
	Image            string         `json:"image"`
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
	CardID     uint32 `json:"cardID"` //服务器返回时候，告知是生成了哪张卡ID
}

// 抢卡结构体
type GetCardData struct {
	GetCardID uint32 `json:"getCardID"`
}
