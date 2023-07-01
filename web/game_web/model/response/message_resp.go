package response

import "game_web/model"

// 服务器返回前端结构体类型,前端通过MsgType字段知道消息是什么类型，做什么处理
type MessageResponse struct {
	MsgType uint32 `json:"msgType"`
	//通用信息
	ChatInfo   ChatResponse   `json:"chatInfo"`
	ErrInfo    ErrResponse    `json:"errInfo"`
	HealthInfo HealthResponse `json:"healthInfo"`
	MsgInfo    MsgResponse    `json:"msgInfo"`
	//游戏信息
	GameStateInfo        GameStateResponse        `json:"gameStateInfo"`
	UserGameInfo         UserGameInfoResponse     `json:"userGameInfo"`
	UseSpecialCardInfo   UseSpecialCardResponse   `json:"useSpecialCardInfo"`
	UseItemInfo          UseItemResponse          `json:"useItemInfo"`
	ScoreRankInfo        ScoreRankResponse        `json:"scoreRankInfo"`
	GameOverInfo         GameOverResponse         `json:"gameOverInfo"`
	GrabCardRoundInfo    GrabCardRoundResponse    `json:"grabCardRoundInfo"`
	SpecialCardRoundInfo SpecialCardRoundResponse `json:"specialCardRoundInfo"`
	//房间信息
	RoomInfo   RoomResponse   `json:"roomInfo"`
	KickerInfo KickerResponse `json:"kickerInfo"`
}

const (
	//通用
	CheckHealthResponseType = 1 << iota
	ChatResponseType        //用户聊天信息
	MsgResponseType
	ErrResponseMsgType
	//房间
	RoomInfoResponseType
	KickerResponseType //被t的人信息
	//游戏
	GameStateResponseType  //游戏状态信息
	UseSpecialCardInfoType //用户使用特殊卡信息
	UseItemResponseType    //用户使用道具信息
	ScoreRankResponseType
	GameOverResponseType
	GrabCardRoundResponseType
	SpecialCardRoundResponseType
)

// 返回的聊天信息（通用）
type ChatResponse struct {
	UserID      uint32            `json:"userID"`
	ChatMsgData model.ChatMsgData `json:"chatMsgData"`
}

type ErrResponse struct {
	Error error `json:"error"`
}

type HealthResponse struct {
}

// 用于给前端返回服务器操作的事情，前端显示给用户出来即可
type MsgResponse struct {
	MsgData string `json:"msgData"` //消息内容
}
