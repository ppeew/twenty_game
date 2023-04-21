package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game_web/global"
	"game_web/model"
	"game_web/proto"
	"game_web/utils"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

var GameData map[uint32]*model.Game = make(map[uint32]*model.Game)

// 游戏主函数
func RunGame(roomID uint32) {
	//游戏初始化阶段
	game := InitGame(roomID)
	//初始化完成，进入游戏主要逻辑
	for i := uint32(0); i < game.GameCount; i++ {
		//循环初始化
		DoFlush(game)
		//发牌阶段
		DoDistributeCard(game)
		//抢卡阶段
		DoListenDistributeCard(game)
		//特殊卡处理阶段
		DoHandleSpecialCard(game)
		//分数计算阶段
		DoScoreCount(game)
	}
	//游戏结束计算排名发奖励阶段
	DoEndGame(game)
	//完成所有环境，退出游戏，回到房间,将当前用户的连接给房间协程
	for userID, info := range game.Users {
		RoomData[roomID].UsersConn[userID] = info.WS
	}
	RoomData[roomID].RecoverChan <- struct{}{}
}

func DoFlush(game *model.Game) {
	game.RandCard = []model.Card{}
	for _, info := range game.Users {
		info.IsGetCard = false
	}
}

func DoEndGame(game *model.Game) {
	type Info struct {
		userID uint32
		score  uint32
	}
	var ranks []Info
	for userID, info := range game.Users {
		ranks = append(ranks, Info{userID: userID, score: info.Score})
	}

	max := uint32(0)
	userID := uint32(0)
	maxIndex := 0
	for i := 0; i < len(ranks); i++ {
		for j := i; i < len(ranks); i++ {
			if ranks[j].score > max {
				max = ranks[i].score
				userID = ranks[i].userID
				maxIndex = j
			}
		}
		temp := Info{
			userID: userID,
			score:  max,
		}
		ranks[maxIndex] = ranks[i]
		ranks[i] = temp
	}
	marshal, err := json.Marshal(ranks)
	if err != nil {
		panic(err)
	}
	BroadcastToAllGameUsers(game, marshal)
}

func DoScoreCount(game *model.Game) {
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

func DoHandleSpecialCard(game *model.Game) {
	//监听用户特殊卡环节,这块要设置超时时间，非一直读取
	select {
	case msg := <-game.CommonChan:
		//正常处理
		switch msg.Type {
		case model.UseSpecialCard:
			//只有这类型的消息才处理
			findCard := false
			for i, card := range game.Users[msg.UserID].SpecialCards {
				if msg.UseSpecialData.SpecialCardID == card.CardID {
					//找到卡，执行
					findCard = true
					game.Users[msg.UserID].SpecialCards = append(game.Users[msg.UserID].SpecialCards[:i], game.Users[msg.UserID].SpecialCards[i+1:]...)
					switch card.Type {
					case model.AddCard:
						data := msg.UseSpecialData.AddCardData
						game.MakeCardID++
						game.Users[msg.UserID].BaseCards = append(game.Users[msg.UserID].BaseCards, model.BaseCard{
							CardID: game.MakeCardID,
							Number: data.NeedNumber,
						})
						ret := fmt.Sprintf("使用增加卡，添加了一张%d的数字卡", data.NeedNumber)
						BroadcastToAllGameUsers(game, []byte(ret))
						break
					case model.DeleteCard:
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
							break
						}
						ret := fmt.Sprintf("使用删除卡，删除了玩家%d一张%d的数字卡", data.TargetUserID, data.CardID)
						BroadcastToAllGameUsers(game, []byte(ret))
					case model.UpdateCard:
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
							break
						}
						ret := fmt.Sprintf("使用更新卡，更新玩家%d一张ID为%d的数字卡为%d", data.TargetUserID, data.CardID, data.UpdateNumber)
						BroadcastToAllGameUsers(game, []byte(ret))
					case model.ChangeCard:
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
								break
							}
						}
						if findUserCard == false {
							utils.SendErrToUser(game.Users[msg.UserID].WS, "[DoHandleSpecialCard]", errors.New("找不到要交换的的卡"))
							break
						}
						for _, info := range game.Users[data.TargetUserID].BaseCards {
							if info.CardID == data.TargetCard {
								findTargetUserCard = true
								targetUserInfo = &info
								break
							}
						}
						if findTargetUserCard == false {
							utils.SendErrToUser(game.Users[msg.UserID].WS, "[DoHandleSpecialCard]", errors.New("找不到对方交换的卡"))
							break
						}
						//都找到了
						temp := userInfo
						userInfo = targetUserInfo
						targetUserInfo = temp
						ret := fmt.Sprintf("玩家%d使用交换卡，用自己%d的卡交换玩家%d一张ID为%d的数字卡", msg.UserID, data.CardID, data.TargetUserID, data.TargetCard)
						BroadcastToAllGameUsers(game, []byte(ret))
					}
					break
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

func DoListenDistributeCard(game *model.Game) {
	//监听用户抢牌环节,这块要设置超时时间，非一直读取
	select {
	case msg := <-game.CommonChan:
		//正常处理
		switch msg.Type {
		//只有这类型的消息才处理
		case model.ListenHandleCard:
			//每一局用户最多只能抢一张卡，检查
			if game.Users[msg.UserID].IsGetCard {
				utils.SendMsgToUser(game.Users[msg.UserID].WS, []byte("一回合只能抢一次噢！"))
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
					utils.SendMsgToUser(game.Users[msg.UserID].WS, []byte("抢到卡了！"))
				} else {
					utils.SendMsgToUser(game.Users[msg.UserID].WS, []byte("没抢到卡~~~"))
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

func DoDistributeCard(game *model.Game) {
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
	for _, info := range game.Users {
		marshal, _ := json.Marshal(game.RandCard)
		err := info.WS.OutChanWrite(marshal)
		if err != nil {
			info.WS.CloseConn()
		}
	}
}

func InitGame(roomID uint32) *model.Game {
	if GameData[roomID] == nil {
		GameData[roomID] = &model.Game{
			RoomID:     roomID,
			Users:      make(map[uint32]*model.UserGameInfo),
			CommonChan: make(chan model.Message, 50),
		}
	}
	game := GameData[roomID]
	room, err := global.GameSrvClient.SearchRoom(context.Background(), &proto.RoomIDInfo{RoomID: game.RoomID})
	if err != nil {
		zap.S().Panic("[RunGame]无法查找到房间信息")
	}
	game.GameCount = room.GameCount
	game.UserNumber = room.UserNumber
	//创建两个协程，用来处理用户聊天消息和用户使用道具，这样才能异步执行
	go ProcessChatMsg(context.TODO(), game)
	go ProcessItemMsg(context.TODO(), game)
	for userID, wsConn := range RoomData[game.RoomID].UsersConn {
		itemsInfo, err := global.GameSrvClient.GetUserItemsInfo(context.Background(), &proto.UserIDInfo{Id: userID})
		if err != nil {
			zap.S().Panic("[RunGame]无法查找到物品信息")
		}
		game.Users[userID] = &model.UserGameInfo{
			BaseCards:    make([]model.BaseCard, 0),
			SpecialCards: make([]model.SpecialCard, 0),
			Items:        itemsInfo.Items,
			WS:           wsConn,
		}
		//对于每个用户开启一个协程，用于读取他的消息到游戏管道（分发消息功能）
		go ReadGameUserMsg(game, userID, wsConn)
	}
	return game
}

func ProcessItemMsg(todo context.Context, game *model.Game) {
	for true {
		select {
		case <-todo.Done():
			//读到主线程停止消息,处理最后的消息,退出
			return
		case item := <-game.ItemChan:
			//处理用户的物品使用,广播所有用户
			msg := []byte(fmt.Sprintf("用户%d对用户%d使用道具%d", item.Item, item.TargetUserID, item.Item))
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
			BroadcastToAllGameUsers(game, msg)
		}
	}
}

func ProcessChatMsg(todo context.Context, game *model.Game) {
	for true {
		select {
		case <-todo.Done():
			//读到主线程停止消息,处理最后的消息,退出
			return
		case chat := <-game.ChatChan:
			//处理用户的聊天消息,广播所有用户
			msg := []byte(fmt.Sprintf("用户%d发送：%s", chat.UserID, string(chat.Data)))
			BroadcastToAllGameUsers(game, msg)
		}
	}
}

func BroadcastToAllGameUsers(game *model.Game, msg []byte) {
	for _, info := range game.Users {
		err := info.WS.OutChanWrite(msg)
		if err != nil {
			info.WS.CloseConn()
		}
	}
}

func ReadGameUserMsg(game *model.Game, userID uint32, wsConn *model.WSConn) {
	for true {
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
		case model.ChatMsgData:
			//聊天信息发到聊天管道
			message.ChatMsg.UserID = userID
			game.ChatChan <- message.ChatMsg
		case model.ItemMsgData:
			message.ItemMsg.UserID = userID
			game.ItemChan <- message.ItemMsg
		default:
			//其他信息是通用信息
			message.UserID = userID
			game.CommonChan <- message
		}
	}
}
