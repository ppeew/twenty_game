package api

import (
	"context"
	"encoding/json"
	"errors"
	"game_web/global"
	"game_web/model"
	"game_web/model/response"
	"game_web/proto"
	"game_web/utils"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

type Game struct {
	RoomID     uint32
	Users      map[uint32]*model.UserGameInfo
	GameCount  uint32
	UserNumber uint32

	CommonChan chan model.Message     //游戏逻辑管道
	ChatChan   chan model.ChatMsgData //聊天管道
	ItemChan   chan model.ItemMsgData //使用物品管道
	HealthChan chan model.Message     //心脏包管道

	MakeCardID uint32       //依次生成卡的id
	RandCard   []model.Card //卡id->卡信息(包含特殊和普通卡)
}

// 游戏号 -> 游戏数据的映射
var GameData map[uint32]*Game = make(map[uint32]*Game)

// 游戏主函数
func RunGame(roomID uint32) {
	//游戏初始化阶段
	game := NewGame(roomID)
	//初始化完成，进入游戏主要逻辑
	for i := uint32(0); i < game.GameCount; i++ {
		//循环初始化
		game.DoFlush()
		time.Sleep(1 * time.Second)
		//发牌阶段
		game.DoDistributeCard()
		//抢卡阶段
		game.DoListenDistributeCard()
		//特殊卡处理阶段
		game.DoHandleSpecialCard()
		//分数计算阶段
		game.DoScoreCount()
	}
	//游戏结束计算排名发奖励阶段
	game.DoEndGame()
	//更改用户为非准备状态，并且房间为等待状态
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: roomID})
	if err != nil {
		zap.S().Infof("err:%s", err)
	}
	room.RoomWait = true
	for _, user := range room.Users {
		user.Ready = false
	}
	for u, _ := range game.Users {
		UsersState[u] = RoomIn
	}
	_, err = global.GameSrvClient.UpdateRoom(context.Background(), room)
	if err != nil {
		zap.S().Info("[RunGame]更新房间失败")
	}
	//完成所有环境，退出游戏，创建房间协程，回到房间协程来
	go startRoomThread(roomID)
}

func NewGame(roomID uint32) *Game {
	if GameData[roomID] == nil {
		GameData[roomID] = &Game{
			RoomID:     roomID,
			Users:      make(map[uint32]*model.UserGameInfo),
			CommonChan: make(chan model.Message, 1024),
			ChatChan:   make(chan model.ChatMsgData, 1024),
			ItemChan:   make(chan model.ItemMsgData, 1024),
			HealthChan: make(chan model.Message, 1024),
		}
	}
	game := GameData[roomID]
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: game.RoomID})
	if err != nil {
		zap.S().Panic("[RunGame]无法查找到房间信息")
	}
	game.GameCount = room.GameCount
	game.UserNumber = room.UserNumber
	//创建三个协程，用来处理用户聊天消息和用户使用道具和心脏包回复，异步执行
	go game.ProcessChatMsg(context.TODO())
	go game.ProcessItemMsg(context.TODO())
	go game.ProcessHealthMsg(context.TODO())
	for u, info := range room.Users {
		itemsInfo, err := global.GameSrvClient.GetUserItemsInfo(context.Background(), &proto.UserIDInfo{Id: info.ID})
		if err != nil {
			zap.S().Panic("[RunGame]无法查找到物品信息")
		}
		game.Users[info.ID] = &model.UserGameInfo{
			BaseCards:    make([]model.BaseCard, 0),
			SpecialCards: make([]model.SpecialCard, 0),
			Items: []uint32{
				itemsInfo.Items[proto.Type_Apple],
				itemsInfo.Items[proto.Type_Banana],
			},
			IsGetCard: false,
			Score:     0,
			WS:        RoomData[roomID].UsersConn[uint32(u)],
		}
		//对于每个用户开启一个协程，用于读取他的消息到游戏管道（分发消息功能）
		go game.ReadGameUserMsg(info.ID)
	}
	BroadcastToAllGameUsers(game, response.RoomMsgResponse{
		MsgType: response.RoomMsgResponseType,
		MsgData: "可以进入游戏了",
	})
	//等待用户的连接,完成，超时都不完成的直接开始游戏
	time.Sleep(6 * time.Second)
	return game
}

func (game *Game) DoFlush() {
	game.RandCard = []model.Card{}
	for _, info := range game.Users {
		info.IsGetCard = false
	}
	resp := CardModelToResponse(game)
	BroadcastToAllGameUsers(game, resp)
}

func (game *Game) DoEndGame() {
	type Info struct {
		UserID uint32
		Score  uint32
	}
	var ranks []Info
	for userID, info := range game.Users {
		ranks = append(ranks, Info{UserID: userID, Score: info.Score})
	}

	max := uint32(0)
	userID := uint32(0)
	maxIndex := 0
	for i := 0; i < len(ranks); i++ {
		for j := i; i < len(ranks); i++ {
			if ranks[j].Score > max {
				max = ranks[i].Score
				userID = ranks[i].UserID
				maxIndex = j
			}
		}
		temp := Info{
			UserID: userID,
			Score:  max,
		}
		ranks[maxIndex] = ranks[i]
		ranks[i] = temp
	}
	BroadcastToAllGameUsers(game, ranks)
}

func (game *Game) DoScoreCount() {
	for _, info := range game.Users {
		sum := uint32(0)
		for i, card := range info.BaseCards {
			sum += card.Number
			//删除这张卡
			info.BaseCards = append(info.BaseCards[:i], info.BaseCards[i+1:]...)
		}
		if sum%12 == 0 {
			info.Score += sum / 12
		}
		//生成多的数字
		game.MakeCardID++
		info.BaseCards = append(info.BaseCards, model.BaseCard{
			CardID: game.MakeCardID,
			Number: sum % 12,
		})
	}
}

func (game *Game) DoHandleSpecialCard() {
	//监听用户特殊卡环节,这块要设置超时时间，非一直读取
	for true {
		select {
		case msg := <-game.CommonChan:
			//正常处理
			switch msg.Type {
			case model.UseSpecialCardMsg:
				//只有这类型的消息才处理
				findCard := false
				for i, card := range game.Users[msg.UserID].SpecialCards {
					if msg.UseSpecialData.SpecialCardID == card.CardID {
						//找到卡，执行
						findCard = true
						game.Users[msg.UserID].SpecialCards = append(game.Users[msg.UserID].SpecialCards[:i], game.Users[msg.UserID].SpecialCards[i+1:]...)
						HandleCard[card.Type](game, msg)
					}
				}
				if findCard == false {
					//找不到卡
					utils.SendErrToUser(game.Users[msg.UserID].WS, "[DoHandleSpecialCard]", errors.New("找不到该特殊卡"))
				}
			default:
				//其他消息不处理,给用户返回超时
				utils.SendErrToUser(game.Users[msg.UserID].WS, "[DoListenDistributeCard]", errors.New("超时信息不处理"))
			}
		case <-time.After(time.Second * 10):
			//超时处理,超时就直接返回了
			return
		}
	}
}

func (game *Game) DoListenDistributeCard() {
	//监听用户抢牌环节,这块要设置超时时间，非一直读取
	for true {
		select {
		case msg := <-game.CommonChan:
			//正常处理
			switch msg.Type {
			//只有这类型的消息才处理
			case model.ListenHandleCardMsg:
				//每一局用户最多只能抢一张卡，检查
				if game.Users[msg.UserID].IsGetCard {
					utils.SendMsgToUser(game.Users[msg.UserID].WS, "一回合只能抢一次噢！")
				} else {
					data := msg.GetCardData
					isOK := false
					for _, card := range game.RandCard {
						if data.GetCardID == card.CardID && !card.HasOwner {
							//按理来说不会空
							//if game.Users[msg.UserID] == nil {
							//	game.Users[msg.UserID] = new(model.UserGameInfo)
							//}
							if card.Type == model.BaseType {
								game.Users[msg.UserID].BaseCards = append(game.Users[msg.UserID].BaseCards, card.BaseCardCardInfo)
							} else if card.Type == model.SpecialType {
								game.Users[msg.UserID].SpecialCards = append(game.Users[msg.UserID].SpecialCards, card.SpecialCardInfo)
							}
							isOK = true
							break
						}
					}
					//发送给用户信息
					if isOK {
						utils.SendMsgToUser(game.Users[msg.UserID].WS, "抢到卡了！")
						resp := CardModelToResponse(game)
						BroadcastToAllGameUsers(game, resp)
					} else {
						utils.SendMsgToUser(game.Users[msg.UserID].WS, "没抢到卡~~~")
					}
				}
			default:
				//其他消息不处理,给用户返回超时
				utils.SendErrToUser(game.Users[msg.UserID].WS, "[DoListenDistributeCard]", errors.New("超时信息不处理"))
			}
		case <-time.After(time.Second * 10):
			//超时处理,超时就直接返回了
			return
		}
	}
}

func (game *Game) DoDistributeCard() {
	//要生成userNumber*2的卡牌，其中包含普通卡和特殊卡,特殊卡数量应该在玩家数量（1/4-1/3）
	needCount := int(game.UserNumber * 2)
	//先生成特殊卡
	rand.Seed(time.Now().Unix())
	special := needCount / (rand.Intn(2) + 3)
	for i := 0; i < special; i++ {
		//生成一张随机的特殊卡
		cardType := 1 << rand.Intn(5)
		game.MakeCardID++
		game.RandCard = append(game.RandCard, model.Card{
			Type: model.SpecialType,
			SpecialCardInfo: model.SpecialCard{
				CardID: game.MakeCardID, //这个字段每张卡必须唯一
				Type:   uint32(cardType),
			},
		})
	}
	//生成普通卡
	needCount -= special
	for i := 0; i < needCount; i++ {
		game.MakeCardID++
		game.RandCard = append(game.RandCard, model.Card{
			Type: model.BaseType,
			BaseCardCardInfo: model.BaseCard{
				CardID: game.MakeCardID,
				Number: uint32(1 + rand.Intn(11)),
			},
		})
	}
	//生成完成,通过websocket发送用户
	resp := CardModelToResponse(game)
	for _, info := range game.Users {
		go utils.SendMsgToUser(info.WS, resp)
	}
}

func (game *Game) ProcessHealthMsg(todo context.Context) {
	for true {
		select {
		case <-todo.Done():
			return
		case msg := <-game.HealthChan:
			utils.SendMsgToUser(game.Users[msg.UserID].WS, response.CheckHealthResponse{
				MsgType: response.CheckHealthResponseType,
				Ok:      true,
			})
		}
	}
}

func (game *Game) ProcessItemMsg(todo context.Context) {
	for true {
		select {
		case <-todo.Done():
			//读到主线程停止消息,处理最后的消息,退出
			return
		case item := <-game.ItemChan:
			items := make([]uint32, 2)
			switch proto.Type(item.Item) {
			case proto.Type_Apple:
				items[proto.Type_Apple] = 1
			case proto.Type_Banana:
				items[proto.Type_Banana] = 1
			}
			isOk, err := global.GameSrvClient.UseItem(context.Background(), &proto.UseItemInfo{
				Id:    item.UserID,
				Items: items,
			})
			if isOk.IsOK == false {
				utils.SendErrToUser(game.Users[item.UserID].WS, "[ProcessItemMsg]", err)
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
			//读到主线程停止消息,处理最后的消息,退出
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

func (game *Game) ReadGameUserMsg(userID uint32) {
	for true {
		wsConn := game.Users[userID].WS
		data, err := wsConn.InChanRead()
		if err != nil {
			//如果读到客户端关闭信息,关闭与客户端的websocket连接
			wsConn.CloseConn()
			continue
		}
		message := model.Message{}
		err = json.Unmarshal(data, &message)
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

func init() {
	HandleCard[model.AddCard] = HandleAddCard
	HandleCard[model.DeleteCard] = HandleDeleteCard
	HandleCard[model.UpdateCard] = HandleUpdateCard
	HandleCard[model.ChangeCard] = HandleChangeCard
}

type HandlerCard func(*Game, model.Message)

var HandleCard map[uint32]HandlerCard = make(map[uint32]HandlerCard)

func HandleChangeCard(game *Game, msg model.Message) {
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
		utils.SendErrToUser(game.Users[msg.UserID].WS, "[DoHandleSpecialCard]", errors.New("找不到要交换的的卡"))
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
		utils.SendErrToUser(game.Users[msg.UserID].WS, "[DoHandleSpecialCard]", errors.New("找不到对方交换的卡"))
		return
	}
	//都找到了
	temp := userInfo
	userInfo = targetUserInfo
	targetUserInfo = temp
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

func HandleUpdateCard(game *Game, msg model.Message) {
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
		utils.SendErrToUser(game.Users[msg.UserID].WS, "[DoHandleSpecialCard]", errors.New("找不到要更新的卡"))
		return
	}
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

func HandleDeleteCard(game *Game, msg model.Message) {
	data := msg.UseSpecialData.DeleteCardData
	findDelCard := false
	for i, card := range game.Users[data.TargetUserID].BaseCards {
		if card.CardID == data.CardID {
			//删除
			findDelCard = true
			game.Users[data.TargetUserID].BaseCards = append(game.Users[data.TargetUserID].BaseCards[:i], game.Users[data.TargetUserID].BaseCards[i+1:]...)
			break
		}
	}
	if findDelCard == false {
		utils.SendErrToUser(game.Users[msg.UserID].WS, "[DoHandleSpecialCard]", errors.New("找不到要删除的卡"))
		return
	}
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

func HandleAddCard(game *Game, msg model.Message) {
	data := msg.UseSpecialData.AddCardData
	game.MakeCardID++
	game.Users[msg.UserID].BaseCards = append(game.Users[msg.UserID].BaseCards, model.BaseCard{
		CardID: game.MakeCardID,
		Number: data.NeedNumber,
	})
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

func BroadcastToAllGameUsers(game *Game, msg interface{}) {
	c := map[string]interface{}{
		"data": msg,
	}
	marshal, _ := json.Marshal(c)
	for _, info := range game.Users {
		err := info.WS.OutChanWrite(marshal)
		if err != nil {
			info.WS.CloseConn()
		}
	}
}

func CardModelToResponse(game *Game) response.GameStateResponse {
	var users []response.UserGameInfoResponse
	for _, info := range game.Users {
		user := response.UserGameInfoResponse{
			BaseCards:    info.BaseCards,
			SpecialCards: info.SpecialCards,
			IsGetCard:    info.IsGetCard,
			Score:        info.Score,
		}
		users = append(users, user)
	}
	resp := response.GameStateResponse{
		MsgType:   response.GameStateResponseType,
		GameCount: game.GameCount,
		Users:     users,
		RandCard:  game.RandCard,
	}
	return resp
}
