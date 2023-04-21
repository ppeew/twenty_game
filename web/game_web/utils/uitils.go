package utils

import (
	"game_web/model"
	"go.uber.org/zap"
	"strings"
)

func StringToBool(str string) bool {
	if strings.ToLower(str) == "true" {
		return true
	}
	return false
}

func SendErrToUser(ws *model.WSConn, handlerFunc string, err error) {
	if err != nil {
		ret := handlerFunc + err.Error()
		err := ws.OutChanWrite([]byte(ret))
		if err != nil {
			zap.S().Infof("ID为%d的用户掉线了", ws.UserID)
		}
	}
}

func SendMsgToUser(ws *model.WSConn, data []byte) {
	err := ws.OutChanWrite(data)
	if err != nil {
		zap.S().Infof("ID为%d的用户掉线了", ws.UserID)
	}
}
