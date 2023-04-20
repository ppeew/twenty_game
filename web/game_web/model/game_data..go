package model

type Game struct {
	Users     map[uint32]UserGameInfo
	GameState uint32
	MsgChan   chan Message
	GameCount uint32
}

// 游戏状态
const (
	StartGame = 1 << iota
	DistributeCard
	ListenDistributeCard
	HandleSpecialCard
	ScoreCount
	EndGame
	Reject
)

type UserGameInfo struct {
	Cards        []uint32
	SpecialCards []uint32
	Items        []uint32
	WS           *WSConn
}
