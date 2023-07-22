package server

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"process_web/global"
	"process_web/my_struct"
	"process_web/my_struct/response"
	game_proto "process_web/proto/game"
	"process_web/utils"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
	"go.uber.org/zap"
)

type GameStruct struct {
	Users        map[uint32]*my_struct.UserGameInfo
	CurrentCount uint32                 //当前是第几回合
	CommonChan   chan my_struct.Message //游戏逻辑管道
	ChatChan     chan my_struct.Message //聊天管道
	ItemChan     chan my_struct.Message //使用物品管道
	HealthChan   chan my_struct.Message //心脏包管道
	wg           sync.WaitGroup
	MakeCardID   uint32            //依次生成卡的id
	RandCard     []*my_struct.Card //卡id->卡信息(包含特殊和普通卡)

	RoomID     uint32
	GameCount  uint32
	UserNumber uint32
	RoomOwner  uint32
	RoomName   string
}

func NewGameData(data *Data) GameStruct {
	users := make(map[uint32]*my_struct.UserGameInfo)
	//查询API用户信息
	for _, userID := range data.users {
		var res utils.UserInfo
		gorequest.New().Get("http://139.159.234.134:8000/user/v1/search").Param("id", strconv.Itoa(int(userID))).
			Retry(5, time.Second, http.StatusInternalServerError).EndStruct(&res)
		users[userID] = &my_struct.UserGameInfo{
			ID:           userID,
			IntoRoomTime: time.Now(),
			Nickname:     res.Nickname,
			Gender:       res.Gender,
			Username:     res.Username,
			Image:        res.Image,

			BaseCards:        make([]*my_struct.BaseCard, 0),
			SpecialCards:     make([]*my_struct.SpecialCard, 0),
			GetBaseCardNum:   0,
			IsGetSpecialCard: false,
			Score:            0,
		}
		zap.S().Infof("[NewGameData]:查询出用户信息%+v", res)
	}

	return GameStruct{
		Users:        users,
		CommonChan:   make(chan my_struct.Message, 1024),
		ChatChan:     make(chan my_struct.Message, 1024),
		ItemChan:     make(chan my_struct.Message, 1024),
		HealthChan:   make(chan my_struct.Message, 1024),
		wg:           sync.WaitGroup{},
		MakeCardID:   0,
		RandCard:     make([]*my_struct.Card, 0),
		RoomID:       data.roomID,
		CurrentCount: 0,
		GameCount:    data.gameCount,
	}
}

func (game *GameStruct) DoFlush() {
	game.RandCard = []*my_struct.Card{}
	for _, info := range game.Users {
		info.GetBaseCardNum = 0
		info.IsGetSpecialCard = false
	}
	game.CurrentCount++
}

func (game *GameStruct) DoDistributeCard() {
	//要生成userNumber+2的卡牌，其中包含普通卡和特殊卡
	needCount := 12
	//特殊卡数量[0,3)
	special := rand.Intn(5)
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
			game.RandCard = append(game.RandCard, &my_struct.Card{
				CardID: game.MakeCardID,
				Type:   my_struct.BaseType,
				BaseCardInfo: my_struct.BaseCard{
					CardID: game.MakeCardID,
					Number: uint32(1 + rand.Intn(9)),
				},
			})
		} else {
			//生成特殊卡
			cardType := 1 << rand.Intn(4)
			game.MakeCardID++
			game.RandCard = append(game.RandCard, &my_struct.Card{
				CardID: game.MakeCardID,
				Type:   my_struct.SpecialType,
				SpecialCardInfo: my_struct.SpecialCard{
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
			userInfo := global.UsersConn[msg.UserID]
			if msg.Type == my_struct.GrabCardMsg {
				data := msg.GetCardData
				isOK := false
				for _, card := range game.RandCard {
					if data.GetCardID == card.CardID && !card.HasOwner {
						if card.Type == my_struct.BaseType {
							if game.Users[msg.UserID].GetBaseCardNum >= 2 {
								global.SendMsgToUser(userInfo, response.MessageResponse{
									MsgType: response.MsgResponseType,
									MsgInfo: &response.MsgResponse{MsgData: "一回合最多抢两次普通卡噢！"},
								})
								return
							}
							game.Users[msg.UserID].BaseCards = append(game.Users[msg.UserID].BaseCards, &card.BaseCardInfo)
							game.Users[msg.UserID].GetBaseCardNum++
						} else if card.Type == my_struct.SpecialType {
							if game.Users[msg.UserID].IsGetSpecialCard {
								global.SendMsgToUser(userInfo, response.MessageResponse{
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
					BroadcastToAllGameUsers(game, CardModelToResponse(game))
				} else {
					global.SendMsgToUser(userInfo, response.MessageResponse{
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
			userInfo := global.UsersConn[msg.UserID]
			if msg.Type == my_struct.UseSpecialCardMsg {
				//只有这类型的消息才处理
				isFind, cardType := game.DropSpecialCard(msg.UserID, msg.UseSpecialData.SpecialCardID)
				if isFind {
					handleFunc[cardType](msg)
				} else {
					global.SendErrToUser(userInfo, "[DoHandleSpecialCard]", errors.New("找不到该特殊卡"))
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
			global.SendMsgToUser(global.UsersConn[id], response.MessageResponse{
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
			info.BaseCards = make([]*my_struct.BaseCard, 0)
			if sum%20 == 0 {
				info.Score += 6
				global.SendMsgToUser(global.UsersConn[id], response.MessageResponse{
					MsgType: response.MsgResponseType,
					MsgInfo: &response.MsgResponse{MsgData: fmt.Sprintf("得分啦！")},
				})
			} else {
				//生成多的数字
				game.MakeCardID++
				info.BaseCards = append(info.BaseCards, &my_struct.BaseCard{
					CardID: game.MakeCardID,
					Number: sum % 20,
				})
				global.SendMsgToUser(global.UsersConn[id], response.MessageResponse{
					MsgType: response.MsgResponseType,
					MsgInfo: &response.MsgResponse{MsgData: fmt.Sprintf("超出20了！")},
				})
			}
		}
	}
	time.Sleep(time.Second)
}

func (game *GameStruct) BackToRoom() {

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
				AddScore:     game.UserNumber - uint32(i), //小心这边变负数
				AddGametimes: 1,
			})
		}
	}()
}

func (game *GameStruct) RunGame() {
	rand.Seed(time.Now().Unix())
	ctx, cancelFunc := context.WithCancel(context.Background())
	//创建三个协程，用来处理用户聊天消息和用户使用道具和心脏包回复，异步执行
	go game.ProcessChatMsg(ctx)
	go game.ProcessItemMsg(ctx)
	go game.ProcessHealthMsg(ctx)
	game.wg.Add(3)
	for _, info := range game.Users {
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
	//游戏初始化阶段
	for i := uint32(0); i < game.GameCount; i++ {
		game.DoFlush()
		BroadcastToAllGameUsers(game, response.MessageResponse{
			MsgType: response.MsgResponseType,
			MsgInfo: &response.MsgResponse{MsgData: "进入抢卡阶段"},
		})
		time.Sleep(time.Second * 2)
		game.DoDistributeCard()
		game.DoListenDistributeCard(10, 14)
		BroadcastToAllGameUsers(game, response.MessageResponse{
			MsgType: response.MsgResponseType,
			MsgInfo: &response.MsgResponse{MsgData: "进入出牌阶段"},
		})
		time.Sleep(time.Second * 2)
		game.DoHandleSpecialCard(14, 18)
		game.DoScoreCount()
	}
	BroadcastToAllGameUsers(game, response.MessageResponse{
		MsgType:      response.GameOverResponseType,
		GameOverInfo: &response.GameOverResponse{},
	})
	cancelFunc() //关闭子协程
	game.wg.Wait()
	//游戏计算排名发奖励阶段
	game.DoEndGame()

}
