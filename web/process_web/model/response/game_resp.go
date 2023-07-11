package response

import (
	"process_web/model"
	"time"
)

// 游戏状态信息(玩家卡牌堆，分数信息)
type GameStateResponse struct {
	GameCount    uint32                 `json:"gameCount"` //游戏总回合数
	GameCurCount uint32                 `json:"gameCurCount"`
	Users        []UserGameInfoResponse `json:"users"`
	RandCard     []*model.Card          `json:"randCard"`
}

type UserGameInfoResponse struct {
	UserID       uint32              `json:"userID"`
	BaseCards    []model.BaseCard    `json:"baseCards"`    //普通卡
	SpecialCards []model.SpecialCard `json:"specialCards"` //特殊卡
	IsGetCard    bool                `json:"isGetCard"`    //当前回合已经抢过卡了
	Score        uint32              `json:"score"`
}

// 使用特殊卡状态信息
type UseSpecialCardResponse struct {
	SpecialCardType uint32               `json:"specialCardType"` //使用特殊卡类型
	UserID          uint32               `json:"userID"`          //使用道具的玩家
	ChangeCardData  model.ChangeCardData `json:"changeCardData"`
	DeleteCardData  model.DeleteCardData `json:"deleteCardData"`
	UpdateCardData  model.UpdateCardData `json:"updateCardData"`
	AddCardData     model.AddCardData    `json:"addCardData"`
}

const (
	AddCard = 1 << iota
	DeleteCard
	UpdateCard
	ChangeCard
)

// 游戏玩家使用道具信息
type UseItemResponse struct {
	ItemMsgData model.ItemMsgData `json:"itemMsgData"`
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
