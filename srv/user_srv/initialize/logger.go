package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		zap.S().Warn("[InitLogger]无法启动日志:%s", err.Error())
	}
	zap.ReplaceGlobals(logger)
}
