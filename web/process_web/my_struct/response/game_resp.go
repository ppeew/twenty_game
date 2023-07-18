package response

import (
	"process_web/my_struct"
	"time"
)

// 游戏状态信息(玩家卡牌堆，分数信息)
type GameStateResponse struct {
	GameCount    uint32                 `json:"gameCount"` //游戏总回合数
	GameCurCount uint32                 `json:"gameCurCount"`
	Users        []UserGameInfoResponse `json:"users"`
	RandCard     []*my_struct.Card      `json:"randCard"`
}

type UserGameInfoResponse struct {
	UserID       uint32                   `json:"userID"`
	BaseCards    []*my_struct.BaseCard    `json:"baseCards"`    //普通卡
	SpecialCards []*my_struct.SpecialCard `json:"specialCards"` //特殊卡
	Score        uint32                   `json:"score"`
	IntoRoomTime time.Time                `json:"-"`
	Nickname     string                   `json:"nickname"`
	Gender       bool                     `json:"gender"`
	Username     string                   `json:"username"`
	Image        string                   `json:"image"`
}

// 使用特殊卡状态信息
type UseSpecialCardResponse struct {
	SpecialCardType uint32                    `json:"specialCardType"` //使用特殊卡类型
	UserID          uint32                    `json:"userID"`          //使用道具的玩家
	ChangeCardData  *my_struct.ChangeCardData `json:"changeCardData,omitempty"`
	DeleteCardData  *my_struct.DeleteCardData `json:"deleteCardData,omitempty"`
	UpdateCardData  *my_struct.UpdateCardData `json:"updateCardData,omitempty"`
	AddCardData     *my_struct.AddCardData    `json:"addCardData,omitempty"`
}

const (
	AddCard = 1 << iota
	DeleteCard
	UpdateCard
	ChangeCard
)

// 游戏玩家使用道具信息
type UseItemResponse struct {
	ItemMsgData my_struct.ItemMsgData `json:"itemMsgData"`
}

type Info struct {
	UserID uint32
	Score  uint32
}

// 游戏结束分数排行信息
type ScoreRankResponse struct {
	Ranks []Info `json:"rank"`
}

// 游戏结束返回消息体
type GameOverResponse struct {
}

// 游戏开始抢卡信息体
type GrabCardRoundResponse struct {
	Duration time.Duration `json:"duration"`
}

// 游戏开始特殊卡处理信息体
type SpecialCardRoundResponse struct {
	Duration time.Duration `json:"duration"`
}
