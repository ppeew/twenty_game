package server

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"process_web/global"
	"process_web/my_struct"
	"process_web/my_struct/response"
	"sort"
)

type HandlerCard func(my_struct.Message)

func NewHandleFunc(game *GameStruct) map[uint32]HandlerCard {
	var HandleCard = make(map[uint32]HandlerCard)
	HandleCard[my_struct.AddCard] = game.HandleAddCard
	HandleCard[my_struct.DeleteCard] = game.HandleDeleteCard
	HandleCard[my_struct.UpdateCard] = game.HandleUpdateCard
	HandleCard[my_struct.ChangeCard] = game.HandleChangeCard
	return HandleCard
}

func (game *GameStruct) HandleAddCard(msg my_struct.Message) {
	data := msg.UseSpecialData.AddCardData
	game.MakeCardID++
	game.Users[msg.UserID].BaseCards = append(game.Users[msg.UserID].BaseCards, &my_struct.BaseCard{
		CardID: game.MakeCardID,
		Number: data.NeedNumber,
	})
	rsp := response.UseSpecialCardResponse{
		SpecialCardType: response.AddCard,
		UserID:          msg.UserID,
		AddCardData: &my_struct.AddCardData{
			NeedNumber: msg.UseSpecialData.AddCardData.NeedNumber,
			CardID:     game.MakeCardID,
		},
	}
	//用户使用增加卡,给该玩家加分
	game.Users[msg.UserID].Score += UseAddScore
	BroadcastToAllGameUsers(game, response.MessageResponse{MsgType: response.UseSpecialCardResponseType, UseSpecialCardInfo: &rsp})
}

func (game *GameStruct) HandleUpdateCard(msg my_struct.Message) {
	//ws := global.UsersConn[msg.UserID]
	value, _ := global.UsersConn.Load(msg.UserID)
	ws := value.(*global.WSConn)
	data := msg.UseSpecialData.UpdateCardData
	findUpdateCard := false
	var findCardNum uint32
	for _, card := range game.Users[data.TargetUserID].BaseCards {
		if card.CardID == data.CardID {
			//更新
			findUpdateCard = true
			findCardNum = card.Number
			card.Number = data.UpdateNumber
		}
	}
	if findUpdateCard == false {
		global.SendErrToUser(ws, "[DoHandleSpecialCard]", errors.New("找不到要更新的卡"))
		return
	}
	game.Users[msg.UserID].Score += 10 - (findCardNum - data.UpdateNumber) + UseUpdateScore
	rsp := response.UseSpecialCardResponse{
		SpecialCardType: response.UpdateCard,
		UserID:          msg.UserID,
		UpdateCardData: &my_struct.UpdateCardData{
			TargetUserID: msg.UseSpecialData.UpdateCardData.TargetUserID,
			CardID:       msg.UseSpecialData.UpdateCardData.CardID,
			UpdateNumber: msg.UseSpecialData.UpdateCardData.UpdateNumber,
		},
	}
	BroadcastToAllGameUsers(game, response.MessageResponse{MsgType: response.UseSpecialCardResponseType, UseSpecialCardInfo: &rsp})

}

func (game *GameStruct) HandleDeleteCard(msg my_struct.Message) {
	//ws := global.UsersConn[msg.UserID]
	value, _ := global.UsersConn.Load(msg.UserID)
	ws := value.(*global.WSConn)
	data := msg.UseSpecialData.DeleteCardData
	findDelCard := false
	if game.Users[data.TargetUserID] == nil {
		global.SendErrToUser(ws, "[HandleDeleteCard]", errors.New("未知的玩家"))
		return
	}
	var delCardNum uint32
	for i, card := range game.Users[data.TargetUserID].BaseCards {
		if card.CardID == data.CardID {
			//删除
			findDelCard = true
			delCardNum = card.Number
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
	// 增加使用炸弹卡玩家得分
	game.Users[msg.UserID].Score += 10 - delCardNum + UseDelScore
	rsp := response.UseSpecialCardResponse{
		SpecialCardType: response.DeleteCard,
		UserID:          msg.UserID,
		DeleteCardData: &my_struct.DeleteCardData{
			TargetUserID: msg.UseSpecialData.DeleteCardData.TargetUserID,
			CardID:       msg.UseSpecialData.DeleteCardData.CardID,
		},
	}
	BroadcastToAllGameUsers(game, response.MessageResponse{MsgType: response.UseSpecialCardResponseType, UseSpecialCardInfo: &rsp})
}

func (game *GameStruct) HandleChangeCard(msg my_struct.Message) {
	//ws := global.UsersConn[msg.UserID]
	value, _ := global.UsersConn.Load(msg.UserID)
	ws := value.(*global.WSConn)
	data := msg.UseSpecialData.ChangeCardData
	//先找到两卡
	findUserCard := false
	findTargetUserCard := false
	var (
		firstChangeNum  uint32
		secondChangeNum uint32
		userInfo        *my_struct.BaseCard
		targetUserInfo  *my_struct.BaseCard
	)
	for _, info := range game.Users[msg.UserID].BaseCards {
		if info.CardID == data.CardID {
			findUserCard = true
			firstChangeNum = info.Number
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
			secondChangeNum = info.Number
			targetUserInfo = info
			break
		}
	}
	if findTargetUserCard == false {
		global.SendErrToUser(ws, "[DoHandleSpecialCard]", errors.New("找不到对方交换的卡"))
		return
	}
	game.Users[msg.UserID].Score += (firstChangeNum+secondChangeNum)/2 + UseChangeScore
	//都找到了
	temp := userInfo
	userInfo = targetUserInfo
	targetUserInfo = temp
	rsp := response.UseSpecialCardResponse{
		SpecialCardType: response.ChangeCard,
		UserID:          msg.UserID,
		ChangeCardData: &my_struct.ChangeCardData{
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
			value, _ := global.UsersConn.Load(msg.UserID)
			ws := value.(*global.WSConn)
			global.SendMsgToUser(ws, response.MessageResponse{
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
		case _ = <-game.ItemChan:
			// TODO 处理item使用逻辑
			//userInfo := global.UsersConn[msg.UserID]
			//items := make([]uint32, 2)
			//switch game_proto.Type(msg.ItemMsgData.Item) {
			//case game_proto.Type_Apple:
			//	items[game_proto.Type_Apple] = 1
			//case game_proto.Type_Banana:
			//	items[game_proto.Type_Banana] = 1
			//}
			//isOk, err := global.GameSrvClient.UseItem(context.Background(), &game_proto.UseItemInfo{
			//	Id:    msg.UserID,
			//	Items: items,
			//})
			//if isOk.IsOK == false {
			//	global.SendErrToUser(userInfo, "[ProcessItemMsg]", err)
			//}
			//处理用户的物品使用,广播所有用户
			//rsp := response.UseItemResponse{
			//	ItemMsgData: my_struct.ItemMsgData{
			//		Item:         msg.ItemMsgData.Item,
			//		TargetUserID: msg.ItemMsgData.TargetUserID,
			//	},
			//}
			//BroadcastToAllGameUsers(game, response.MessageResponse{MsgType: response.UseItemResponseType, UseItemInfo: &rsp})
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

func (game *GameStruct) ForUserIntoRoom(ctx context.Context) {
	if global.IntoRoomCHAN[game.RoomID] == nil {
		global.IntoRoomCHAN[game.RoomID] = make(chan uint32)
	}
	for true {
		select {
		case <-ctx.Done():
			game.wg.Done()
			return
		case _ = <-global.IntoRoomCHAN[game.RoomID]:
			global.IntoRoomRspCHAN[game.RoomID] <- false
		}
	}
}

// 读取用户信息协程
func (game *GameStruct) ReadGameUserMsg(ctx context.Context, userID uint32) {
	for true {
		value, _ := global.UsersConn.Load(userID)
		ws := value.(*global.WSConn)
		zap.S().Info("[ReadGameUserMsg]:等待ws信息中")
		select {
		case <-ctx.Done():
			zap.S().Info("[ReadGameUserMsg]:收到退出信号")
			game.wg.Done()
			return
		case message := <-ws.InChanRead():
			switch message.Type {
			case my_struct.ChatMsg:
				//聊天信息发到聊天管道
				message.UserID = userID
				game.ChatChan <- message
			case my_struct.ItemMsg:
				//物品信息发到物品管道
				message.UserID = userID
				game.ItemChan <- message
			case my_struct.GetGameMsg:
				global.SendMsgToUser(ws, CardModelToResponse(game))
			case my_struct.GetState:
				//获取状态
				global.SendMsgToUser(ws, response.MessageResponse{
					MsgType:      response.GetStateResponseType,
					GetStateInfo: &response.GetStateResponse{State: 1},
				})
			//case model.CheckHealthMsg:
			//	//心脏包
			//	message.UserID = ShopID
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
		value, _ := global.UsersConn.Load(userID)
		ws := value.(*global.WSConn)
		err := ws.OutChanWrite(msg)
		if err != nil {
			//global.UsersConn[ShopID].CloseConn()
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
			Score:        info.Score,
			IntoRoomTime: info.IntoRoomTime,
			Nickname:     info.Nickname,
			Gender:       info.Gender,
			Username:     info.Username,
			Image:        info.Image,
		}
		users = append(users, userGameInfoResponse)
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].IntoRoomTime.Before(users[j].IntoRoomTime)
	})
	info := response.GameStateResponse{
		GameCount:    game.GameCount,
		GameCurCount: game.CurrentCount,
		Users:        users,
		RandCard:     game.RandCard,
	}
	return response.MessageResponse{MsgType: response.GameStateResponseType, GameStateInfo: &info}
}
