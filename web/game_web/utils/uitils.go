package utils

import (
	"game_web/model"
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
			ws.CloseConn()
		}
	}
}

func SendMsgToUser(ws *model.WSConn, data []byte) {
	err := ws.OutChanWrite(data)
	if err != nil {
		ws.CloseConn()
	}
}
