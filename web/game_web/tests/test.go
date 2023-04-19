package tests

import (
	"encoding/json"
	"fmt"
	"game_web/model"
	"github.com/gorilla/websocket"
)

func main() {
	dial, _, err := websocket.DefaultDialer.Dial("ws://localhost:8083/ws", nil)
	if err != nil {
		panic(err)
	}
	defer dial.Close()
	mes := model.Message{
		UserID:     0,
		Type:       model.DeleteRoom,
		DeleteData: model.DeleteData{RoomID: 1234},
		UpdateData: model.UpdateData{},
		RoomData:   model.RoomData{},
		ReadyState: model.ReadyState{},
		BeginGame:  model.BeginGame{},
	}
	marshal, _ := json.Marshal(mes)
	err = dial.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		panic(err)
	}
	_, p, err := dial.ReadMessage()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(p))
}
