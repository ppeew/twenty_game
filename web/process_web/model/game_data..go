package model

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
	BaseCards        []*BaseCard    `json:"baseCards,omitempty"`    //普通卡
	SpecialCards     []*SpecialCard `json:"specialCards,omitempty"` //特殊卡
	IsGetCard        bool           `json:"isGetCard,omitempty"`    //当前回合已经抢过普通卡了
	IsGetSpecialCard bool           `json:"isGetSpecialCard"`       //当前回合已经抢过特殊卡了
	Score            uint32         `json:"score,omitempty"`
	IntoRoomTime     time.Time      `json:"-"`
}
