package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"game_web/model"
	"game_web/model/response"
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
		err2 := ws.OutChanWrite(marshal)
		if err2 != nil {
			zap.S().Infof("ID为%d的用户掉线了", ws.UserID)
		}
	}
}

func SendMsgToUser(ws *model.WSConn, data interface{}) {
	c := map[string]interface{}{
		"data": data,
	}
	zap.S().Infof("[SendMsgToUser]:正在向用户%d发送信息,消息为:%v", ws.UserID, data)
	marshal, _ := json.Marshal(c)
	err := ws.OutChanWrite(marshal)
	if err != nil {
		zap.S().Infof("ID为%d的用户掉线了", ws.UserID)
	}
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
