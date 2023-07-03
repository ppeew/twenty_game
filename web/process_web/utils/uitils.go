package utils

import (
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
