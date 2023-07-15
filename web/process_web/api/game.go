package api

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"process_web/global"
	"process_web/model"
	"process_web/model/response"
	game_proto "process_web/proto/game"
	"sort"
	"sync"
	"time"
)

type GameStruct struct {
	Users        map[uint32]*model.UserGameInfo
	CurrentCount uint32             //当前是第几回合
	CommonChan   chan model.Message //游戏逻辑管道
	ChatChan     chan model.Message //聊天管道
	ItemChan     chan model.Message //使用物品管道
	HealthChan   chan model.Message //心脏包管道
	MakeCardID   uint32             //依次生成卡的id
	RandCard     []*model.Card      //卡id->卡信息(包含特殊和普通卡)
	exitCancel   context.CancelFunc //负责退出
	wg           sync.WaitGroup     //等待其他协程退出

	GameData GameData
}

type GameData struct {
	RoomID     uint32
	GameCount  uint32
	UserNumber uint32
	RoomOwner  uint32
	Users      map[uint32]response.UserData
	RoomName   string
}

// RunGame 游戏主函数
func RunGame(data GameData) {
	//游戏初始化阶段
	game := NewGame(data)
	for i := uint32(0); i < game.GameData.GameCount; i++ {
		//循环初始化
		game.DoFlush()
		//zap.S().Info("游戏[DoFlush]完成")
		BroadcastToAllGameUsers(game, response.MessageResponse{
			MsgType: response.MsgResponseType,
			MsgInfo: &response.MsgResponse{MsgData: "进入抢卡阶段"},
		})
		time.Sleep(time.Second * 2)
		//zap.S().Info("游戏[BroadcastToAllGameUsers]完成")
		//发牌阶段
		game.DoDistributeCard()
		//zap.S().Info("游戏[DoDistributeCard]完成")
		//抢卡阶段
		game.DoListenDistributeCard(6, 8)
		//zap.S().Info("游戏[DoListenDistributeCard]完成")
		BroadcastToAllGameUsers(game, response.MessageResponse{
			MsgType: response.MsgResponseType,
			MsgInfo: &response.MsgResponse{MsgData: "进入出牌阶段"},
		})
		time.Sleep(time.Second * 2)
		//特殊卡处理阶段
		game.DoHandleSpecialCard(8, 18)
		//zap.S().Info("游戏[DoHandleSpecialCard]完成")
		//分数计算阶段
		game.DoScoreCount()
		//zap.S().Info("游戏[DoScoreCount]完成")
	}
	// 回到房间
	game.BackToRoom()
	//游戏计算排名发奖励阶段
	game.DoEndGame()
	//zap.S().Info("游戏[RunGame]完成")
}

func NewGame(data GameData) *GameStruct {
	rand.Seed(time.Now().Unix())
	ctx, cancelFunc := context.WithCancel(context.Background())
	game := &GameStruct{
		Users:        make(map[uint32]*model.UserGameInfo),
		CurrentCount: 0,
		CommonChan:   make(chan model.Message, 1024),
		ChatChan:     make(chan model.Message, 1024),
		ItemChan:     make(chan model.Message, 1024),
		HealthChan:   make(chan model.Message, 1024),
		exitCancel:   cancelFunc,
		wg:           sync.WaitGroup{},
		GameData:     data,
	}
	//创建三个协程，用来处理用户聊天消息和用户使用道具和心脏包回复，异步执行
	go game.ProcessChatMsg(ctx)
	go game.ProcessItemMsg(ctx)
	go game.ProcessHealthMsg(ctx)
	game.wg.Add(3)
	for _, info := range game.GameData.Users {
		//zap.S().Infof("[NewGame]初始化用户%d", info.ID)
		game.Users[info.ID] = &model.UserGameInfo{
			BaseCards:    make([]*model.BaseCard, 0),
			SpecialCards: make([]*model.SpecialCard, 0),
			IsGetCard:    false,
			Score:        0,
			IntoRoomTime: info.IntoRoomTime,
		}
		//对于每个用户开启一个协程，用于读取他的消息到游戏管道（分发消息功能）
		go game.ReadGameUserMsg(ctx, info.ID)
		game.wg.Add(1)
	}
	//等待用户页面初始化完成
	BroadcastToAllGameUsers(game, response.MessageResponse{
		MsgType: response.MsgResponseType,
		MsgInfo: &response.MsgResponse{MsgData: "游戏2秒后开始！"},
	})
	time.Sleep(time.Second * 2)
	return game
}

func (game *GameStruct) DoFlush() {
	game.RandCard = []*model.Card{}
	for _, info := range game.Users {
		info.IsGetCard = false
		info.IsGetSpecialCard = false
	}
	game.CurrentCount++
}

func (game *GameStruct) DoDistributeCard() {
	//要生成userNumber+2的卡牌，其中包含普通卡和特殊卡
	needCount := 6
	//特殊卡数量[0,1]
	special := rand.Intn(3)
	cards := make([]int, needCount)
	used := make([]bool, needCount)
	for i := needCount - 1; i >= needCount-special; i-- {
		cards[i] = 1
	}
	for i := 0; i < needCount; i++ {
		index := rand.Intn(needCount)
		for used[index] {
			index = (index + 1) % needCount
		}
		//持续过程，直到找到第一个没使用的
		used[index] = true
		if cards[index] == 0 {
			//生成普通卡
			game.MakeCardID++
			game.RandCard = append(game.RandCard, &model.Card{
				CardID: game.MakeCardID,
				Type:   model.BaseType,
				BaseCardInfo: model.BaseCard{
					CardID: game.MakeCardID,
					Number: uint32(1 + rand.Intn(9)),
				},
			})
		} else {
			//生成特殊卡
			cardType := 1 << rand.Intn(4)
			game.MakeCardID++
			game.RandCard = append(game.RandCard, &model.Card{
				CardID: game.MakeCardID,
				Type:   model.SpecialType,
				SpecialCardInfo: model.SpecialCard{
					CardID: game.MakeCardID, //这个字段每张卡必须唯一
					Type:   uint32(cardType),
				},
			})
		}
	}
	//生成完成,通过websocket发送用户
	BroadcastToAllGameUsers(game, CardModelToResponse(game))
}

func (game *GameStruct) DoListenDistributeCard(min, max int) {
	duration := time.Duration(rand.Intn(max-min)+min) * time.Second
	//监听用户抢牌环节,这块要设置超时时间，非一直读取
	after := time.After(duration)
	BroadcastToAllGameUsers(game, response.MessageResponse{
		MsgType: response.GrabCardRoundResponseType,
		GrabCardRoundInfo: &response.GrabCardRoundResponse{
			Duration: duration,
		},
	})
	for true {
		select {
		case msg := <-game.CommonChan:
			userInfo := UsersConn[msg.UserID]
			if msg.Type == model.GrabCardMsg {
				data := msg.GetCardData
				isOK := false
				for _, card := range game.RandCard {
					if data.GetCardID == card.CardID && !card.HasOwner {
						if card.Type == model.BaseType {
							if game.Users[msg.UserID].IsGetCard {
								SendMsgToUser(userInfo, response.MessageResponse{
									MsgType: response.MsgResponseType,
									MsgInfo: &response.MsgResponse{MsgData: "一回合最多抢一次普通卡噢！"},
								})
								return
							}
							game.Users[msg.UserID].BaseCards = append(game.Users[msg.UserID].BaseCards, &card.BaseCardInfo)
							game.Users[msg.UserID].IsGetCard = true
						} else if card.Type == model.SpecialType {
							if game.Users[msg.UserID].IsGetSpecialCard {
								SendMsgToUser(userInfo, response.MessageResponse{
									MsgType: response.MsgResponseType,
									MsgInfo: &response.MsgResponse{MsgData: "一回合最多抢一次特殊卡噢！"},
								})
								return
							}
							game.Users[msg.UserID].SpecialCards = append(game.Users[msg.UserID].SpecialCards, &card.SpecialCardInfo)
							game.Users[msg.UserID].IsGetSpecialCard = true
						}
						card.HasOwner = true //range的value是值拷贝！！！
						isOK = true
						break
					}
				}
				//发送给用户信息
				if isOK {
					SendMsgToUser(userInfo, response.MessageResponse{
						MsgType: response.MsgResponseType,
						MsgInfo: &response.MsgResponse{MsgData: "抢到卡了！"},
					})
					resp := CardModelToResponse(game)
					BroadcastToAllGameUsers(game, resp)
				} else {
					SendMsgToUser(userInfo, response.MessageResponse{
						MsgType: response.MsgResponseType,
						MsgInfo: &response.MsgResponse{MsgData: "没抢到卡~~~"},
					})
				}
			}
		case <-after:
			//超时处理,超时就直接返回了
			return
		}
	}
}

func (game *GameStruct) DoHandleSpecialCard(min, max int) {
	//如果所有用户都没有特殊卡，回合加快
	var speed = true
	for _, data := range game.Users {
		if len(data.SpecialCards) > 0 {
			speed = false
			break
		}
	}

	duration := time.Duration(rand.Intn(max-min)+min) * time.Second
	if speed {
		duration = time.Second * 3
	}
	BroadcastToAllGameUsers(game, response.MessageResponse{
		MsgType: response.SpecialCardRoundResponseType,
		SpecialCardRoundInfo: &response.SpecialCardRoundResponse{
			Duration: duration,
		},
	})
	//监听用户特殊卡环节,这块要设置超时时间，非一直读取
	handleFunc := NewHandleFunc(game)
	after := time.After(duration)
	for true {
		select {
		case msg := <-game.CommonChan:
			userInfo := UsersConn[msg.UserID]
			if msg.Type == model.UseSpecialCardMsg {
				//只有这类型的消息才处理
				isFind, cardType := game.DropSpecialCard(msg.UserID, msg.UseSpecialData.SpecialCardID)
				if isFind {
					handleFunc[cardType](msg)
				} else {
					SendErrToUser(userInfo, "[DoHandleSpecialCard]", errors.New("找不到该特殊卡"))
				}
			}
		case <-after:
			//超时处理,超时就直接返回了
			return
		}
	}
}

func (game *GameStruct) DoScoreCount() {
	for id, info := range game.Users {
		//首先清理用户普通卡（要求：普通卡不能大于6张，大于6张则删除最先进来的卡）
		total := len(info.BaseCards)
		if total > 6 {
			info.BaseCards = info.BaseCards[total-6:]
			SendMsgToUser(UsersConn[id], response.MessageResponse{
				MsgType: response.MsgResponseType,
				MsgInfo: &response.MsgResponse{MsgData: fmt.Sprintf("卡没了噢爆了！")},
			})
		}
		//处理分数
		sum := uint32(0)
		for _, card := range info.BaseCards {
			sum += card.Number
		}
		if sum/20 == 1 {
			info.BaseCards = make([]*model.BaseCard, 0)
			if sum%20 == 0 {
				info.Score += 6
				SendMsgToUser(UsersConn[id], response.MessageResponse{
					MsgType: response.MsgResponseType,
					MsgInfo: &response.MsgResponse{MsgData: fmt.Sprintf("得分啦！")},
				})
			} else {
				//生成多的数字
				game.MakeCardID++
				info.BaseCards = append(info.BaseCards, &model.BaseCard{
					CardID: game.MakeCardID,
					Number: sum % 20,
				})
				SendMsgToUser(UsersConn[id], response.MessageResponse{
					MsgType: response.MsgResponseType,
					MsgInfo: &response.MsgResponse{MsgData: fmt.Sprintf("超出12了！")},
				})
			}
		}
	}
	time.Sleep(time.Second)
}

func (game *GameStruct) BackToRoom() {
	BroadcastToAllGameUsers(game, response.MessageResponse{
		MsgType:      response.GameOverResponseType,
		GameOverInfo: &response.GameOverResponse{},
	})
	game.exitCancel() //关闭子协程
	game.wg.Wait()    //等待全部子协程关闭

	//完成所有环境，退出游戏协程，创建房间协程
	go startRoomThread(RoomData{
		RoomID:        game.GameData.RoomID,
		MaxUserNumber: game.GameData.UserNumber,
		GameCount:     game.GameData.GameCount,
		UserNumber:    game.GameData.UserNumber,
		RoomOwner:     game.GameData.RoomOwner,
		Users:         game.GameData.Users,
		RoomName:      game.GameData.RoomName,
		RoomWait:      true,
	})
}

func (game *GameStruct) DoEndGame() {
	var ranks []response.Info
	for userID, info := range game.Users {
		ranks = append(ranks, response.Info{UserID: userID, Score: info.Score})
	}
	sort.Slice(ranks, func(i, j int) bool {
		if ranks[i].Score == ranks[j].Score {
			return ranks[i].UserID > ranks[j].UserID
		}
		return ranks[i].Score > ranks[j].Score
	})
	BroadcastToAllGameUsers(game, response.MessageResponse{
		MsgType: response.ScoreRankResponseType,
		ScoreRankInfo: &response.ScoreRankResponse{
			Ranks: ranks,
		},
	})
	// 分别给第一名和其他名次玩家发放奖励(第一名发放钻石，第一第二名，发放物品，全部玩家名发放金币),计算排行榜
	go func() {
		for i, u := range ranks {
			global.GameSrvClient.UpdateRanks(context.Background(), &game_proto.UpdateRanksInfo{
				UserID:       u.UserID,
				AddScore:     game.GameData.UserNumber - uint32(i),
				AddGametimes: 1,
			})
		}
	}()
}
