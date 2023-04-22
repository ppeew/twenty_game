package response

import "game_web/model"

// websocket返回结构体类型,前端通过MsgType字段知道消息是什么类型，做什么处理
const (
	GameStateResponseType = 1 << iota
	UseSpecialCardInfoType
	UseItemResponseType
	ChatResponseType
)

// 游戏状态信息(玩家卡牌堆，分数信息)
type GameStateResponse struct {
	MsgType   uint32                 `json:"msgType"`
	GameCount uint32                 `json:"gameCount"`
	Users     []UserGameInfoResponse `json:"users"`
	RandCard  []model.Card           `json:"randCard"`
}

type UserGameInfoResponse struct {
	BaseCards    []model.BaseCard    `json:"baseCards"`    //普通卡
	SpecialCards []model.SpecialCard `json:"specialCards"` //特殊卡
	IsGetCard    bool                `json:"isGetCard"`    //当前回合已经抢过卡了
	Score        uint32              `json:"score"`
}

// 使用特殊卡状态信息
type UseSpecialCardResponse struct {
	MsgType         uint32               `json:"msgType"`
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
	MsgType     uint32            `json:"msgType"`
	ItemMsgData model.ItemMsgData `json:"itemMsgData"`
}

// 游戏玩家聊天信心
type ChatResponse struct {
	MsgType     uint32            `json:"msgType"`
	ChatMsgData model.ChatMsgData `json:"chatMsgData"`
}
