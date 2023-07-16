package server

import (
	"context"
	"errors"
	"process_web/global"
	"process_web/model"
	"process_web/model/response"
	game_proto "process_web/proto/game"
	"sort"
)

type HandlerCard func(model.Message)

func NewHandleFunc(game *GameStruct) map[uint32]HandlerCard {
	var HandleCard = make(map[uint32]HandlerCard)
	HandleCard[model.AddCard] = game.HandleAddCard
	HandleCard[model.DeleteCard] = game.HandleDeleteCard
	HandleCard[model.UpdateCard] = game.HandleUpdateCard
	HandleCard[model.ChangeCard] = game.HandleChangeCard
	return HandleCard
}

func (game *GameStruct) HandleAddCard(msg model.Message) {
	data := msg.UseSpecialData.AddCardData
	game.MakeCardID++
	game.Users[msg.UserID].BaseCards = append(game.Users[msg.UserID].BaseCards, &model.BaseCard{
		CardID: game.MakeCardID,
		Number: data.NeedNumber,
	})
	rsp := response.UseSpecialCardResponse{
		SpecialCardType: response.AddCard,
		UserID:          msg.UserID,
		AddCardData: &model.AddCardData{
			NeedNumber: msg.UseSpecialData.AddCardData.NeedNumber,
			CardID:     game.MakeCardID,
		},
	}
	BroadcastToAllGameUsers(game, response.MessageResponse{MsgType: response.UseSpecialCardResponseType, UseSpecialCardInfo: &rsp})
}

func (game *GameStruct) HandleUpdateCard(msg model.Message) {
	ws := global.UsersConn[msg.UserID]
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
		global.SendErrToUser(ws, "[DoHandleSpecialCard]", errors.New("找不到要更新的卡"))
		return
	}
	//game.DropSpecialCard(msg.UserID, msg.UseSpecialData.SpecialCardID)
	rsp := response.UseSpecialCardResponse{
		SpecialCardType: response.UpdateCard,
		UserID:          msg.UserID,
		UpdateCardData: &model.UpdateCardData{
			TargetUserID: msg.UseSpecialData.UpdateCardData.TargetUserID,
			CardID:       msg.UseSpecialData.UpdateCardData.CardID,
			UpdateNumber: msg.UseSpecialData.UpdateCardData.UpdateNumber,
		},
	}
	BroadcastToAllGameUsers(game, response.MessageResponse{MsgType: response.UseSpecialCardResponseType, UseSpecialCardInfo: &rsp})

}

func (game *GameStruct) HandleDeleteCard(msg model.Message) {
	ws := global.UsersConn[msg.UserID]
	data := msg.UseSpecialData.DeleteCardData
	findDelCard := false
	//zap.S().Infof("[HandleDeleteCard]:玩家包括%v", game.Users)
	//zap.S().Infof("[HandleDeleteCard]:被删除卡的玩家是%d", data.TargetUserID)
	if game.Users[data.TargetUserID] == nil {
		global.SendErrToUser(global.UsersConn[msg.UserID], "[HandleDeleteCard]", errors.New("未知的玩家"))
		return
	}
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
		global.SendErrToUser(ws, "[DoHandleSpecialCard]", errors.New("找不到要删除的卡"))
		return
	}
	//game.DropSpecialCard(msg.UserID, msg.UseSpecialData.SpecialCardID)
	rsp := response.UseSpecialCardResponse{
		SpecialCardType: response.DeleteCard,
		UserID:          msg.UserID,
		DeleteCardData: &model.DeleteCardData{
			TargetUserID: msg.UseSpecialData.DeleteCardData.TargetUserID,
			CardID:       msg.UseSpecialData.DeleteCardData.CardID,
		},
	}
	BroadcastToAllGameUsers(game, response.MessageResponse{MsgType: response.UseSpecialCardResponseType, UseSpecialCardInfo: &rsp})
}

func (game *GameStruct) HandleChangeCard(msg model.Message) {
	ws := global.UsersConn[msg.UserID]
	data := msg.UseSpecialData.ChangeCardData
	//先找到两卡
	findUserCard := false
	var userInfo *model.BaseCard
	findTargetUserCard := false
	var targetUserInfo *model.BaseCard
	for _, info := range game.Users[msg.UserID].BaseCards {
		if info.CardID == data.CardID {
			findUserCard = true
			userInfo = info
			break
		}
	}
	if findUserCard == false {
		global.SendErrToUser(ws, "[DoHandleSpecialCard]", errors.New("找不到要交换的的卡"))
		return
	}
	for _, info := range game.Users[data.TargetUserID].BaseCards {
		if info.CardID == data.TargetCard {
			findTargetUserCard = true
			targetUserInfo = info
			break
		}
	}
	if findTargetUserCard == false {
		global.SendErrToUser(ws, "[DoHandleSpecialCard]", errors.New("找不到对方交换的卡"))
		return
	}
	//都找到了
	temp := userInfo
	userInfo = targetUserInfo
	targetUserInfo = temp
	rsp := response.UseSpecialCardResponse{
		SpecialCardType: response.ChangeCard,
		UserID:          msg.UserID,
		ChangeCardData: &model.ChangeCardData{
			CardID:       msg.UseSpecialData.ChangeCardData.CardID,
			TargetUserID: msg.UseSpecialData.ChangeCardData.TargetUserID,
			TargetCard:   msg.UseSpecialData.ChangeCardData.TargetCard,
		},
	}
	BroadcastToAllGameUsers(game, response.MessageResponse{MsgType: response.UseSpecialCardResponseType, UseSpecialCardInfo: &rsp})
}

func (game *GameStruct) ProcessHealthMsg(todo context.Context) {
	for true {
		select {
		case <-todo.Done():
			game.wg.Done()
			return
		case msg := <-game.HealthChan:
			global.SendMsgToUser(global.UsersConn[msg.UserID], response.MessageResponse{
				MsgType:         response.CheckHealthType,
				HealthCheckInfo: &response.HealthCheck{},
			})
		}
	}
}

func (game *GameStruct) ProcessItemMsg(todo context.Context) {
	for true {
		select {
		case <-todo.Done():
			//zap.S().Info("[ProcessItemMsg]退出")
			game.wg.Done()
			return
		case msg := <-game.ItemChan:
			userInfo := global.UsersConn[msg.UserID]
			items := make([]uint32, 2)
			switch game_proto.Type(msg.ItemMsgData.Item) {
			case game_proto.Type_Apple:
				items[game_proto.Type_Apple] = 1
			case game_proto.Type_Banana:
				items[game_proto.Type_Banana] = 1
			}
			isOk, err := global.GameSrvClient.UseItem(context.Background(), &game_proto.UseItemInfo{
				Id:    msg.UserID,
				Items: items,
			})
			if isOk.IsOK == false {
				global.SendErrToUser(userInfo, "[ProcessItemMsg]", err)
			}
			//处理用户的物品使用,广播所有用户
			rsp := response.UseItemResponse{
				ItemMsgData: model.ItemMsgData{
					Item:         msg.ItemMsgData.Item,
					TargetUserID: msg.ItemMsgData.TargetUserID,
				},
			}
			BroadcastToAllGameUsers(game, response.MessageResponse{MsgType: response.UseItemResponseType, UseItemInfo: &rsp})
		}
	}
}

func (game *GameStruct) ProcessChatMsg(todo context.Context) {
	for true {
		select {
		case <-todo.Done():
			//zap.S().Info("[ProcessChatMsg]退出")
			game.wg.Done()
			return
		case msg := <-game.ChatChan:
			//处理用户的聊天消息,广播所有用户
			rsp := response.ChatResponse{
				UserID:      msg.UserID,
				ChatMsgData: msg.ChatMsgData.Data,
			}
			BroadcastToAllGameUsers(game, response.MessageResponse{MsgType: response.ChatResponseType, ChatInfo: &rsp})
		}
	}
}

// 读取用户信息协程
func (game *GameStruct) ReadGameUserMsg(ctx context.Context, userID uint32) {
	for true {
		select {
		case <-ctx.Done():
			game.wg.Done()
			return
		case message := <-global.UsersConn[userID].InChanRead():
			switch message.Type {
			case model.ChatMsg:
				//聊天信息发到聊天管道
				message.UserID = userID
				game.ChatChan <- message
			case model.ItemMsg:
				//物品信息发到物品管道
				message.UserID = userID
				game.ItemChan <- message
			//case model.CheckHealthMsg:
			//	//心脏包
			//	message.UserID = userID
			//	game.HealthChan <- message
			default:
				//其他信息是通用信息
				message.UserID = userID
				game.CommonChan <- message
			}
		}
	}
}

func (game *GameStruct) DropSpecialCard(userID uint32, specialID uint32) (bool, uint32) {
	isFind := false
	var cardType uint32
	user := game.Users[userID]
	for index, specialCard := range user.SpecialCards {
		if specialCard.CardID == specialID {
			if index+1 >= len(user.SpecialCards) {
				user.SpecialCards = user.SpecialCards[:index]
			} else {
				user.SpecialCards = append(user.SpecialCards[:index], user.SpecialCards[index+1:]...)
			}
			isFind = true
			cardType = specialCard.Type
			break
		}
	}
	return isFind, cardType
}

func BroadcastToAllGameUsers(game *GameStruct, msg response.MessageResponse) {
	for userID := range game.Users {
		//zap.S().Infof("[BroadcastToAllGameUsers]:正在向用户%d发送信息,消息为:%v", userID, msg)
		err := global.UsersConn[userID].OutChanWrite(msg)
		if err != nil {
			//zap.S().Infof("[BroadcastToAllGameUsers]:%d用户关闭了连接", userID)
			//global.UsersConn[userID].CloseConn()
		}
	}
}

func CardModelToResponse(game *GameStruct) response.MessageResponse {
	var users []response.UserGameInfoResponse
	for userID, info := range game.Users {
		userGameInfoResponse := response.UserGameInfoResponse{
			UserID:       userID,
			BaseCards:    info.BaseCards,
			SpecialCards: info.SpecialCards,
			//IsGetCard:    info.GetBaseCardNum,
			Score:        info.Score,
			IntoRoomTime: info.IntoRoomTime,
		}
		users = append(users, userGameInfoResponse)
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].IntoRoomTime.Before(users[j].IntoRoomTime)
	})
	info := response.GameStateResponse{
		GameCount:    game.GameData.GameCount,
		GameCurCount: game.CurrentCount,
		Users:        users,
		RandCard:     game.RandCard,
	}
	return response.MessageResponse{MsgType: response.GameStateResponseType, GameStateInfo: &info}
}
