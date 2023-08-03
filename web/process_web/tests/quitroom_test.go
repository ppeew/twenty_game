package tests

//
//import (
//	"fmt"
//	"math/rand"
//	"process_web/my_struct/response"
//	"process_web/server"
//	"testing"
//
//	"go.uber.org/zap"
//)
//
//func TestQuit(t *testing.T) {
//	//输入
//	us := make(map[uint32]response.UserData, 2)
//	us[4] = response.UserData{
//		ShopID:    4,
//		Ready: false,
//	}
//	us[5] = response.UserData{
//		ShopID:    5,
//		Ready: false,
//	}
//	roomInfo := server.RoomStruct{
//		RoomData: server.RoomData{
//			RoomID:        2,
//			MaxUserNumber: 2,
//			GameCount:     3,
//			UserNumber:    2,
//			RoomOwner:     4,
//			RoomWait:      false,
//			Users:         us,
//			RoomName:      "test666",
//		},
//	}
//
//	//处理 玩家退出
//	delete(roomInfo.RoomData.Users, 5)
//	roomInfo.RoomData.UserNumber--
//	zap.S().Infof("[QuitRoom]:%d", roomInfo.RoomData.UserNumber)
//	if roomInfo.RoomData.UserNumber == 0 {
//		//没人了，销毁房间
//		fmt.Println("销毁房间了！")
//		return
//	}
//	if 5 == roomInfo.RoomData.RoomOwner {
//		//是房主,转移房间
//		num := rand.Intn(int(roomInfo.RoomData.UserNumber))
//		for _, data := range roomInfo.RoomData.Users {
//			if num <= 0 {
//				roomInfo.RoomData.RoomOwner = data.ShopID
//				break
//			}
//			num--
//		}
//	}
//
//	fmt.Printf("%+v\n", roomInfo.RoomData)
//
//	//处理 玩家退出
//	delete(roomInfo.RoomData.Users, 4)
//	roomInfo.RoomData.UserNumber--
//	zap.S().Infof("[QuitRoom]:%d", roomInfo.RoomData.UserNumber)
//	if roomInfo.RoomData.UserNumber == 0 {
//		//没人了，销毁房间
//		fmt.Println("销毁房间了！")
//		return
//	}
//	if 4 == roomInfo.RoomData.RoomOwner {
//		//是房主,转移房间
//		num := rand.Intn(int(roomInfo.RoomData.UserNumber))
//		for _, data := range roomInfo.RoomData.Users {
//			if num <= 0 {
//				roomInfo.RoomData.RoomOwner = data.ShopID
//				break
//			}
//			num--
//		}
//	}
//
//	//输出
//	fmt.Printf("%+v\n", roomInfo.RoomData)
//}
