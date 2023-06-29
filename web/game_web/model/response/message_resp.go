package response

import "game_web/model"

// websocket返回结构体类型,前端通过MsgType字段知道消息是什么类型，做什么处理
const (
	//通用
	CheckHealthResponseType = 1 << iota
	ChatResponseType        //用户聊天信息
	//房间
	RoomInfoResponseType
	KickerResponseType //被t的人信息
	RoomMsgResponseType
	//游戏
	GameStateResponseType  //游戏状态信息
	UseSpecialCardInfoType //用户使用特殊卡信息
	UseItemResponseType    //用户使用道具信息
	ScoreRankResponseType
	GameOverResponseType
	BeginListenDistributeCard
	BeginHandleSpecialCard
	//错误类型信息
	ErrMsg
)

// 返回的聊天信息（通用）
type ChatResponse struct {
	MsgType     uint32            `json:"msgType"`
	UserID      uint32            `json:"userID"`
	ChatMsgData model.ChatMsgData `json:"chatMsgData"`
}

type ErrData struct {
	MsgType uint32 `json:"msgType"`
	Error   error  `json:"error"`
}

type CheckHealthResponse struct {
	MsgType uint32 `json:"msgType"`
}
