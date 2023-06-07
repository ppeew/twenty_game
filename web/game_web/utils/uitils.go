package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game_web/api"
	"game_web/global"
	"game_web/model"
	"game_web/model/response"
	game_proto "game_web/proto/game"
	"game_web/proto/user"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
)

func StringToBool(str string) bool {
	if strings.ToLower(str) == "true" {
		return true
	}
	return false
}

func SendErrToUser(ws *model.WSConn, handlerFunc string, error error) {
	if error != nil {
		errMsg := response.ErrData{MsgType: response.ErrMsg, Error: errors.New(fmt.Sprintf("[%s]:%s", handlerFunc, error))}
		c := map[string]interface{}{
			"data": errMsg,
		}
		marshal, _ := json.Marshal(c)
		_ = ws.OutChanWrite(marshal)
	}
}

func SendMsgToUser(ws *model.WSConn, data interface{}) {
	c := map[string]interface{}{
		"data": data,
	}
	zap.S().Infof("[SendMsgToUser]:正在向用户发送信息,消息为:%v", data)
	marshal, _ := json.Marshal(c)
	_ = ws.OutChanWrite(marshal)
}

func CheckGoRoutines() {
	go func() {
		for true {
			select {
			case <-time.After(time.Second):
				zap.S().Infof("协程数量->%d", runtime.NumGoroutine())
			}
		}
	}()
}

func ReleaseResource() {
	//释放房间
	for roomID, _ := range api.CHAN {
		_, err := global.GameSrvClient.DeleteRoom(context.Background(), &game_proto.RoomIDInfo{RoomID: roomID})
		if err != nil {
			zap.S().Infof("[ReleaseResource]:关闭%d房间失败:%s", roomID, err.Error())
		}
	}
	//释放用户状态
	for userID, _ := range api.UsersState {
		_, err := global.UserSrvClient.UpdateUserState(context.Background(), &user.UpdateUserStateInfo{Id: userID, State: api.OutSide})
		if err != nil {
			zap.S().Infof("[ReleaseResource]:%s", err)
		}
	}
}
