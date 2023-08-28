package response

// 服务器返回前端结构体类型,前端通过MsgType字段知道消息是什么类型，做什么处理
type MessageResponse struct {
	MsgType uint32 `json:"msgType"`
	//通用信息
	HealthCheckInfo *HealthCheck      `json:"healthCheckInfo,omitempty"` //服务器发送的健康检查包\
	ChatInfo        *ChatResponse     `json:"chatInfo,omitempty"`
	ErrInfo         *ErrResponse      `json:"errInfo,omitempty"`
	MsgInfo         *MsgResponse      `json:"msgInfo,omitempty"`
	GetStateInfo    *GetStateResponse `json:"getStateInfo,omitempty"`
	//游戏信息
	GameStateInfo        *GameStateResponse        `json:"gameStateInfo,omitempty"`
	UserGameInfo         *UserGameInfoResponse     `json:"userGameInfo,omitempty"`
	UseSpecialCardInfo   *UseSpecialCardResponse   `json:"useSpecialCardInfo,omitempty"`
	UseItemInfo          *UseItemResponse          `json:"useItemInfo,omitempty"`
	ScoreRankInfo        *ScoreRankResponse        `json:"scoreRankInfo,omitempty"`
	GameOverInfo         *GameOverResponse         `json:"gameOverInfo,omitempty"`
	GrabCardRoundInfo    *GrabCardRoundResponse    `json:"grabCardRoundInfo,omitempty"`
	SpecialCardRoundInfo *SpecialCardRoundResponse `json:"specialCardRoundInfo,omitempty"`
	//房间信息
	RoomInfo      *RoomResponse   `json:"roomInfo,omitempty"`
	KickerInfo    *KickerResponse `json:"kickerInfo,omitempty"`
	BeginGameInfo *BeginGameData  `json:"beginGameInfo,omitempty"`
}

const (
	//通用
	CheckHealthType      = 100 + iota //心脏包消息
	ChatResponseType                  //用户聊天信息
	MsgResponseType                   //服务器处理完成的消息（打印给用户看即可）
	ErrResponseMsgType                //错误返回消息
	GetStateResponseType              //获取重连状态信息
)

const (
	//房间
	RoomInfoResponseType  = 200 + iota //房间信息
	KickerResponseType                 //T的人信息
	BeginGameResponseType              //开始游戏提醒消息
)

const (
	//游戏
	GameStateResponseType        = 300 + iota //游戏状态信息
	UseSpecialCardResponseType                //用户使用特殊卡信息
	UseItemResponseType                       //用户使用道具信息
	ScoreRankResponseType                     //游戏结束排名信息
	GameOverResponseType                      //游戏结束信息
	GrabCardRoundResponseType                 //用户抢卡回合时间信息
	SpecialCardRoundResponseType              //特殊卡回合时间信息
)

type GetStateResponse struct {
	State int `json:"state"` //0:在房间 1:在单人对战游戏
}

// 返回的聊天信息（通用）
type ChatResponse struct {
	UserID      uint32 `json:"userID"`
	ChatMsgData string `json:"chatMsgData"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

type HealthCheck struct {
}

// 用于给前端返回服务器操作的事情，前端显示给用户出来即可
type MsgResponse struct {
	StateType int    `json:"stateType"`
	MsgData   string `json:"msgData"` //消息内容
}
