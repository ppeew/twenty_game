package utils

import (
	"encoding/json"
	"file_web/model"
	"strings"

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
		c := map[string]interface{}{
			"err": handlerFunc + error.Error(),
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
	marshal, _ := json.Marshal(c)
	err := ws.OutChanWrite(marshal)
	if err != nil {
		zap.S().Infof("ID为%d的用户掉线了", ws.UserID)
	}
}
