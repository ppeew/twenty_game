package server

import (
	"context"
	"process_web/global"
	"process_web/model/response"
	"process_web/proto/game"
	"time"
)

// 房间主函数
func StartRoomThread(data RoomData) {
	ctx, cancel := context.WithCancel(context.Background())
	room := NewRoom(data)
	dealFunc := NewDealFunc(room)
	//读取房间内的管道 (正常来说，用户进入房间但是还没建立socket，此时连接为nil,该读取协程会关闭，当用户游戏结束，连接不为nil)
	for _, userData := range room.RoomData.Users {
		go room.ReadRoomUserMsg(ctx, userData.ID)
	}
	//用于用户连接
	go room.ForUserConn(ctx)
	//用于用户进房
	go room.ForUserIntoRoom(ctx)
	//定时检查房间用户是否占用房间不退出（看socket是否断开了）
	go room.CheckClientHealth(ctx)
	//定时发送到redis，更新房间列表信息，为大厅外查询更新数据
	go room.UpdateRedisRoom(ctx)
	for {
		select {
		case msg := <-room.MsgChan:
			//zap.S().Infof("[StartRoomThread]]:%+v", msg)
			if dealFunc[msg.Type] != nil {
				dealFunc[msg.Type](msg)
			}
		case msg := <-room.ExitChan:
			// 停止信号，关闭主函数及相关子协程，优雅退出
			cancel()
			room.wg.Wait()
			//zap.S().Info("[StartRoomThread]]:其他协程已关闭")
			if msg == RoomQuit {
				global.GameSrvClient.DeleteRoom(context.Background(), &game.RoomIDInfo{RoomID: room.RoomData.RoomID})
				global.GameSrvClient.DelRoomServer(context.Background(), &game.RoomIDInfo{RoomID: room.RoomData.RoomID})
				return
			} else if msg == GameStart {
				room.RoomData.RoomWait = false
				var users []*game.RoomUser
				for _, data := range room.RoomData.Users {
					users = append(users, &game.RoomUser{
						ID:    data.ID,
						Ready: data.Ready,
					})
				}
				global.GameSrvClient.SetGlobalRoom(context.Background(), &game.RoomInfo{
					RoomID:        room.RoomData.RoomID,
					MaxUserNumber: room.RoomData.MaxUserNumber,
					GameCount:     room.RoomData.GameCount,
					UserNumber:    room.RoomData.UserNumber,
					RoomOwner:     room.RoomData.RoomOwner,
					RoomWait:      room.RoomData.RoomWait,
					Users:         users,
					RoomName:      room.RoomData.RoomName,
				})
				go RunGame(GameData{
					RoomID:     room.RoomData.RoomID,
					GameCount:  room.RoomData.GameCount,
					UserNumber: room.RoomData.MaxUserNumber,
					RoomOwner:  room.RoomData.RoomOwner,
					Users:      room.RoomData.Users,
					RoomName:   room.RoomData.RoomName,
				})
				return
			}
		}
	}
}

// RunGame 游戏主函数
func RunGame(data GameData) {
	//游戏初始化阶段
	game := NewGame(data)
	for i := uint32(0); i < game.GameData.GameCount; i++ {
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
	// 回到房间
	game.BackToRoom()
	//游戏计算排名发奖励阶段
	game.DoEndGame()
}
