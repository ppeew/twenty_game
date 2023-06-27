package api

import (
	"context"
	"encoding/json"
	"errors"
	"game_web/global"
	"game_web/model"
	"game_web/model/response"
	game_proto "game_web/proto/game"
	"game_web/utils"
	"math/rand"
	"sort"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Game struct {
	RoomID       uint32
	Users        map[uint32]*model.UserGameInfo
	GameCount    uint32 //游戏总回合数
	CurrentCount uint32 //当前是第几回合
	UserNumber   uint32
	CommonChan   chan model.Message     //游戏逻辑管道
	ChatChan     chan model.ChatMsgData //聊天管道
	ItemChan     chan model.ItemMsgData //使用物品管道
	HealthChan   chan model.Message     //心脏包管道
	MakeCardID   uint32                 //依次生成卡的id
	RandCard     []model.Card           //卡id->卡信息(包含特殊和普通卡)

	exitCancel context.CancelFunc //负责退出
	wg         sync.WaitGroup     //等待其他协程退出
}

func NewGame(roomID uint32) *Game {
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &game_proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		zap.S().Panic("[RunGame]无法查找到房间信息")
	}
	rand.Seed(time.Now().Unix())
	ctx, cancelFunc := context.WithCancel(context.Background())
	game := &Game{
		RoomID:       roomID,
		Users:        make(map[uint32]*model.UserGameInfo),
		GameCount:    room.GameCount,
		CurrentCount: 0,
		UserNumber:   room.UserNumber,
		CommonChan:   make(chan model.Message, 1024),
		ChatChan:     make(chan model.ChatMsgData, 1024),
		ItemChan:     make(chan model.ItemMsgData, 1024),
		HealthChan:   make(chan model.Message, 1024),
		exitCancel:   cancelFunc,
		wg:           sync.WaitGroup{},
	}
	//创建三个协程，用来处理用户聊天消息和用户使用道具和心脏包回复，异步执行
	go game.ProcessChatMsg(ctx)
	go game.ProcessItemMsg(ctx)
	go game.ProcessHealthMsg(ctx)
	game.wg.Add(3)
	for _, roomUser := range room.Users {
		zap.S().Infof("[NewGame]初始化用户%d", roomUser.ID)
		game.Users[roomUser.ID] = &model.UserGameInfo{
			BaseCards:    make([]model.BaseCard, 0),
			SpecialCards: make([]model.SpecialCard, 0),
			IsGetCard:    false,
			Score:        0,
		}
		//对于每个用户开启一个协程，用于读取他的消息到游戏管道（分发消息功能）
		go game.ReadGameUserMsg(ctx, roomUser.ID)
		game.wg.Add(1)
	}
	//等待用户页面初始化完成
	time.Sleep(time.Second)
	return game
}

// RunGame 游戏主函数
func RunGame(roomID uint32) {
	//游戏初始化阶段
	game := NewGame(roomID)
	zap.S().Info("游戏[NewGame]完成")
	for i := uint32(0); i < game.GameCount; i++ {
		//循环初始化
		game.DoFlush()
		//zap.S().Info("游戏[DoFlush]完成")
		BroadcastToAllGameUsers(game, CardModelToResponse(game))
		//zap.S().Info("游戏[BroadcastToAllGameUsers]完成")
		//发牌阶段
		game.DoDistributeCard()
		//zap.S().Info("游戏[DoDistributeCard]完成")
		//抢卡阶段
		//time.Sleep(time.Second * 2)
		game.DoListenDistributeCard(6, 8)
		//zap.S().Info("游戏[DoListenDistributeCard]完成")
		time.Sleep(time.Second * 2)
		//特殊卡处理阶段
		game.DoHandleSpecialCard(6, 18)
		//zap.S().Info("游戏[DoHandleSpecialCard]完成")
		//分数计算阶段
		game.DoScoreCount()
		//zap.S().Info("游戏[DoScoreCount]完成")
	}
	// 回到房间
	game.BackToRoom()
	//游戏计算排名发奖励阶段
	game.DoEndGame()
	zap.S().Info("游戏[RunGame]完成")
}

func (game *Game) DoFlush() {
	game.RandCard = []model.Card{}
	for _, info := range game.Users {
		info.IsGetCard = false
	}
	game.CurrentCount++
}

func (game *Game) DoEndGame() {
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
	BroadcastToAllGameUsers(game, response.ScoreRankResponse{
		MsgType: response.ScoreRankResponseType,
		Ranks:   ranks,
	})
	// 分别给第一名和其他名次玩家发放奖励(第一名发放钻石，第一第二名，发放物品，全部玩家名发放金币) TODO
}

func (game *Game) BackToRoom() {
	//更改用户为非准备状态，并且房间为等待状态
	_, err := global.GameSrvClient.BackRoom(context.Background(), &game_proto.RoomIDInfo{RoomID: game.RoomID})
	if err != nil {
		zap.S().Infof("err:%s", err)
	}
	//for u := range game.Users {
	//	_, _ = global.UserSrvClient.UpdateUserState(context.Background(), &user_proto.UpdateUserStateInfo{Id: u, State: RoomIn})
	//}
	BroadcastToAllGameUsers(game, response.GameOverResponse{MsgType: response.GameOverResponseType})
	//完成所有环境，退出游戏协程，创建房间协程，回到房间协程来
	game.exitCancel() //关闭子协程
	//zap.S().Info("[BackToRoom]等待其他子协程关闭")
	game.wg.Wait() //等待全部子协程关闭
	//zap.S().Info("[BackToRoom]其他子协程已关闭")
	go startRoomThread(game.RoomID)
}

func (game *Game) DoScoreCount() {
	for _, info := range game.Users {
		//首先清理用户普通卡（要求：普通卡不能大于6张，大于6张则删除最先进来的卡）
		total := len(info.BaseCards)
		if total > 6 {
			info.BaseCards = info.BaseCards[total-6:]
		}
		//处理分数
		sum := uint32(0)
		for _, card := range info.BaseCards {
			sum += card.Number
		}
		if sum/12 == 1 {
			info.BaseCards = []model.BaseCard{}
			if sum%12 == 0 {
				info.Score += 5
			} else {
				//生成多的数字
				game.MakeCardID++
				info.BaseCards = append(info.BaseCards, model.BaseCard{
					CardID: game.MakeCardID,
					Number: sum % 12,
				})
			}
		}
	}
}

func (game *Game) DoHandleSpecialCard(min, max int) {
	duration := time.Duration(rand.Intn(max-min)+min) * time.Second
	BroadcastToAllGameUsers(game, response.BeginHandleSpecialCardResponse{
		MsgType:  response.BeginHandleSpecialCard,
		Duration: duration,
	})
	//监听用户特殊卡环节,这块要设置超时时间，非一直读取
	handleFunc := NewHandleFunc(game)
	after := time.After(duration)
	for true {
		select {
		case msg := <-game.CommonChan:
			userInfo := UsersState[msg.UserID]
			switch msg.Type {
			case model.UseSpecialCardMsg:
				//只有这类型的消息才处理
				findCard := false
				for i, card := range game.Users[msg.UserID].SpecialCards {
					if msg.UseSpecialData.SpecialCardID == card.CardID {
						//找到卡，执行
						findCard = true
						game.Users[msg.UserID].SpecialCards = append(game.Users[msg.UserID].SpecialCards[:i], game.Users[msg.UserID].SpecialCards[i+1:]...)
						handleFunc[card.Type](msg)
					}
				}
				if findCard == false {
					//找不到卡
					utils.SendErrToUser(userInfo, "[DoHandleSpecialCard]", errors.New("找不到该特殊卡"))
				}
			default:
				//其他消息不处理,给用户返回超时
				utils.SendErrToUser(userInfo, "[DoListenDistributeCard]", errors.New("其他信息不处理"))
			}
		case <-after:
			//超时处理,超时就直接返回了
			return
		}
	}
}

func (game *Game) DoListenDistributeCard(min, max int) {
	duration := time.Duration(rand.Intn(max-min)+min) * time.Second
	BroadcastToAllGameUsers(game, response.BeginListenDistributeCardResponse{
		MsgType:  response.BeginListenDistributeCard,
		Duration: duration,
	})
	//监听用户抢牌环节,这块要设置超时时间，非一直读取
	after := time.After(duration)
	for true {
		select {
		case msg := <-game.CommonChan:
			userInfo := UsersState[msg.UserID]
			//zap.S().Infof("[DoListenDistributeCard]:%+v", msg)
			switch msg.Type {
			//只有这类型的消息才处理
			case model.ListenHandleCardMsg:
				//每一局用户最多只能抢一张卡，检查
				if game.Users[msg.UserID].IsGetCard {
					utils.SendMsgToUser(userInfo, "一回合只能抢一次噢！")
				} else {
					data := msg.GetCardData
					isOK := false
					//zap.S().Infof("[DoListenDistributeCard]:%+v", data)
					for _, card := range game.RandCard {
						//zap.S().Infof("[DoListenDistributeCard]:%v", card.HasOwner)
						if data.GetCardID == card.CardID && !card.HasOwner {
							if card.Type == model.BaseType {
								game.Users[msg.UserID].BaseCards = append(game.Users[msg.UserID].BaseCards, card.BaseCardInfo)
							} else if card.Type == model.SpecialType {
								game.Users[msg.UserID].SpecialCards = append(game.Users[msg.UserID].SpecialCards, card.SpecialCardInfo)
							}
							card.HasOwner = true
							game.Users[msg.UserID].IsGetCard = true
							isOK = true
							break
						}
					}
					//发送给用户信息
					if isOK {
						utils.SendMsgToUser(userInfo, "抢到卡了！")
						resp := CardModelToResponse(game)
						BroadcastToAllGameUsers(game, resp)
					} else {
						utils.SendMsgToUser(userInfo, "没抢到卡~~~")
					}
				}
			}
		case <-after:
			//超时处理,超时就直接返回了
			return
		}
	}
}

func (game *Game) DoDistributeCard() {
	//要生成userNumber+2的卡牌，其中包含普通卡和特殊卡
	needCount := int(game.UserNumber + 2)
	//特殊卡数量[0,1]
	special := rand.Intn(2)
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
			game.RandCard = append(game.RandCard, model.Card{
				CardID: game.MakeCardID,
				Type:   model.BaseType,
				BaseCardInfo: model.BaseCard{
					CardID: game.MakeCardID,
					Number: uint32(1 + rand.Intn(11)),
				},
			})
		} else {
			//生成特殊卡
			cardType := 1 << rand.Intn(4)
			game.MakeCardID++
			game.RandCard = append(game.RandCard, model.Card{
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
	resp := CardModelToResponse(game)
	BroadcastToAllGameUsers(game, resp)
}

func (game *Game) ProcessHealthMsg(todo context.Context) {
	for true {
		select {
		case <-todo.Done():
			//zap.S().Info("[ProcessHealthMsg]退出")
			game.wg.Done()
			return
		case msg := <-game.HealthChan:
			utils.SendMsgToUser(UsersState[msg.UserID], response.CheckHealthResponse{
				MsgType: response.CheckHealthResponseType,
			})
		}
	}
}

func (game *Game) ProcessItemMsg(todo context.Context) {
	for true {
		select {
		case <-todo.Done():
			//zap.S().Info("[ProcessItemMsg]退出")
			game.wg.Done()
			return
		case item := <-game.ItemChan:
			userInfo := UsersState[item.UserID]
			items := make([]uint32, 2)
			switch game_proto.Type(item.Item) {
			case game_proto.Type_Apple:
				items[game_proto.Type_Apple] = 1
			case game_proto.Type_Banana:
				items[game_proto.Type_Banana] = 1
			}
			isOk, err := global.GameSrvClient.UseItem(context.Background(), &game_proto.UseItemInfo{
				Id:    item.UserID,
				Items: items,
			})
			if isOk.IsOK == false {
				utils.SendErrToUser(userInfo, "[ProcessItemMsg]", err)
			}
			//处理用户的物品使用,广播所有用户
			rsp := response.UseItemResponse{
				MsgType: response.UseItemResponseType,
				ItemMsgData: model.ItemMsgData{
					UserID:       item.UserID,
					Item:         item.Item,
					TargetUserID: item.TargetUserID,
				},
			}
			BroadcastToAllGameUsers(game, rsp)
		}
	}
}

func (game *Game) ProcessChatMsg(todo context.Context) {
	for true {
		select {
		case <-todo.Done():
			//zap.S().Info("[ProcessChatMsg]退出")
			game.wg.Done()
			return
		case chat := <-game.ChatChan:
			//处理用户的聊天消息,广播所有用户
			rsp := response.ChatResponse{
				UserID:  chat.UserID,
				MsgType: response.ChatResponseType,
				ChatMsgData: model.ChatMsgData{
					UserID: chat.UserID,
					Data:   chat.Data,
				},
			}
			BroadcastToAllGameUsers(game, rsp)
		}
	}
}

func (game *Game) ReadGameUserMsg(ctx context.Context, userID uint32) {
	for true {
		select {
		case <-ctx.Done():
			zap.S().Info("[ReadGameUserMsg]退出")
			game.wg.Done()
			return
		case data := <-UsersState[userID].InChanRead():
			message := model.Message{}
			err := json.Unmarshal(data, &message)
			if err != nil {
				//客户端发过来数据有误
				zap.S().Info("客户端发送数据有误:", data)
				continue
			}
			switch message.Type {
			case model.ChatMsg:
				//聊天信息发到聊天管道
				message.ChatMsgData.UserID = userID
				game.ChatChan <- message.ChatMsgData
			case model.ItemMsg:
				//物品信息发到物品管道
				message.ItemMsgData.UserID = userID
				game.ItemChan <- message.ItemMsgData
			case model.CheckHealthMsg:
				//心脏包
				message.UserID = userID
				game.HealthChan <- message
			default:
				//其他信息是通用信息
				message.UserID = userID
				game.CommonChan <- message
			}
		}
	}
}

func (game *Game) HandleChangeCard(msg model.Message) {
	ws := UsersState[msg.UserID]
	data := msg.UseSpecialData.ChangeCardData
	//先找到两卡
	findUserCard := false
	var userInfo *model.BaseCard
	findTargetUserCard := false
	var targetUserInfo *model.BaseCard
	for _, info := range game.Users[msg.UserID].BaseCards {
		if info.CardID == data.CardID {
			findUserCard = true
			userInfo = &info
			return
		}
	}
	if findUserCard == false {
		utils.SendErrToUser(ws, "[DoHandleSpecialCard]", errors.New("找不到要交换的的卡"))
		return
	}
	for _, info := range game.Users[data.TargetUserID].BaseCards {
		if info.CardID == data.TargetCard {
			findTargetUserCard = true
			targetUserInfo = &info
			return
		}
	}
	if findTargetUserCard == false {
		utils.SendErrToUser(ws, "[DoHandleSpecialCard]", errors.New("找不到对方交换的卡"))
		return
	}
	//都找到了
	temp := userInfo
	userInfo = targetUserInfo
	targetUserInfo = temp
	game.DropSpecialCard(msg.UserID, msg.UseSpecialData.SpecialCardID)
	rsp := response.UseSpecialCardResponse{
		MsgType:         response.UseSpecialCardInfoType,
		SpecialCardType: response.ChangeCard,
		UserID:          msg.UserID,
		ChangeCardData: model.ChangeCardData{
			CardID:       msg.UseSpecialData.ChangeCardData.CardID,
			TargetUserID: msg.UseSpecialData.ChangeCardData.TargetUserID,
			TargetCard:   msg.UseSpecialData.ChangeCardData.TargetCard,
		},
	}
	BroadcastToAllGameUsers(game, rsp)

}

func (game *Game) HandleUpdateCard(msg model.Message) {
	ws := UsersState[msg.UserID]
	data := msg.UseSpecialData.UpdateCardData
	findUpdateCard := false
	for _, card := range game.Users[data.TargetUserID].BaseCards {
		if card.CardID == data.CardID {
			//更新
			findUpdateCard = true
			card.Number = data.UpdateNumber
		}
	}
	if findUpdateCard == false {
		utils.SendErrToUser(ws, "[DoHandleSpecialCard]", errors.New("找不到要更新的卡"))
		return
	}
	game.DropSpecialCard(msg.UserID, msg.UseSpecialData.SpecialCardID)
	rsp := response.UseSpecialCardResponse{
		MsgType:         response.UseSpecialCardInfoType,
		SpecialCardType: response.UpdateCard,
		UserID:          msg.UserID,
		UpdateCardData: model.UpdateCardData{
			TargetUserID: msg.UseSpecialData.UpdateCardData.TargetUserID,
			CardID:       msg.UseSpecialData.UpdateCardData.CardID,
			UpdateNumber: msg.UseSpecialData.UpdateCardData.UpdateNumber,
		},
	}
	BroadcastToAllGameUsers(game, rsp)

}

func (game *Game) HandleDeleteCard(msg model.Message) {
	ws := UsersState[msg.UserID]
	data := msg.UseSpecialData.DeleteCardData
	findDelCard := false
	for i, card := range game.Users[data.TargetUserID].BaseCards {
		if card.CardID == data.CardID {
			//删除
			findDelCard = true
			if i+1 >= len(game.Users[data.TargetUserID].BaseCards) {
				game.Users[data.TargetUserID].BaseCards = game.Users[data.TargetUserID].BaseCards[:i]
			} else {
				game.Users[data.TargetUserID].BaseCards = append(game.Users[data.TargetUserID].BaseCards[:i], game.Users[data.TargetUserID].BaseCards[i+1:]...)
			}
			break
		}
	}
	if findDelCard == false {
		utils.SendErrToUser(ws, "[DoHandleSpecialCard]", errors.New("找不到要删除的卡"))
		return
	}
	game.DropSpecialCard(msg.UserID, msg.UseSpecialData.SpecialCardID)
	rsp := response.UseSpecialCardResponse{
		MsgType:         response.UseSpecialCardInfoType,
		SpecialCardType: response.DeleteCard,
		UserID:          msg.UserID,
		DeleteCardData: model.DeleteCardData{
			TargetUserID: msg.UseSpecialData.DeleteCardData.TargetUserID,
			CardID:       msg.UseSpecialData.DeleteCardData.CardID,
		},
	}
	BroadcastToAllGameUsers(game, rsp)
}

func (game *Game) HandleAddCard(msg model.Message) {
	data := msg.UseSpecialData.AddCardData
	game.MakeCardID++
	game.Users[msg.UserID].BaseCards = append(game.Users[msg.UserID].BaseCards, model.BaseCard{
		CardID: game.MakeCardID,
		Number: data.NeedNumber,
	})
	game.DropSpecialCard(msg.UserID, msg.UseSpecialData.SpecialCardID)
	rsp := response.UseSpecialCardResponse{
		MsgType:         response.UseSpecialCardInfoType,
		SpecialCardType: response.AddCard,
		UserID:          msg.UserID,
		AddCardData: model.AddCardData{
			NeedNumber: msg.UseSpecialData.AddCardData.NeedNumber,
		},
	}
	BroadcastToAllGameUsers(game, rsp)
}

func (game *Game) DropSpecialCard(userID uint32, specialID uint32) {
	for u, info := range game.Users {
		if u == userID {
			for index, specialCard := range info.SpecialCards {
				if specialCard.CardID == specialID {
					if index+1 >= len(info.SpecialCards) {
						info.SpecialCards = info.SpecialCards[:index]
					} else {
						info.SpecialCards = append(info.SpecialCards[:index], info.SpecialCards[index+1:]...)
					}
					break
				}
			}
			break
		}
	}
}

func BroadcastToAllGameUsers(game *Game, msg interface{}) {
	c := map[string]interface{}{
		"data": msg,
	}
	marshal, _ := json.Marshal(c)
	for userID := range game.Users {
		zap.S().Infof("[BroadcastToAllGameUsers]:正在向用户%d发送信息,消息为:%v", userID, msg)
		err := UsersState[userID].OutChanWrite(marshal)
		if err != nil {
			zap.S().Infof("[BroadcastToAllGameUsers]:%d用户关闭了连接", userID)
			UsersState[userID].CloseConn()
		}
	}
}

func CardModelToResponse(game *Game) response.GameStateResponse {
	var users []response.UserGameInfoResponse
	for userID, info := range game.Users {
		userGameInfoResponse := response.UserGameInfoResponse{
			UserID:       userID,
			BaseCards:    info.BaseCards,
			SpecialCards: info.SpecialCards,
			IsGetCard:    info.IsGetCard,
			Score:        info.Score,
		}
		users = append(users, userGameInfoResponse)
	}
	resp := response.GameStateResponse{
		MsgType:      response.GameStateResponseType,
		GameCount:    game.GameCount,
		GameCurCount: game.CurrentCount,
		Users:        users,
		RandCard:     game.RandCard,
	}
	return resp
}

type HandlerCard func(model.Message)

func NewHandleFunc(game *Game) map[uint32]HandlerCard {
	var HandleCard = make(map[uint32]HandlerCard)
	HandleCard[model.AddCard] = game.HandleAddCard
	HandleCard[model.DeleteCard] = game.HandleDeleteCard
	HandleCard[model.UpdateCard] = game.HandleUpdateCard
	HandleCard[model.ChangeCard] = game.HandleChangeCard
	return HandleCard
}
