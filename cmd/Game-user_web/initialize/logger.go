package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"math/rand"
	"os"
	"strings"
	"time"
	"user_web/global"
)

const (
	logDir = "log"
)

var logFile = new(os.File)

// InitLogger 每次启动都创建个新日志
func InitLogger() {
	logName := getLoggerName(global.ServerConfig.ConsulInfo.Name)
	zap.S().Infof("[InitLogger] logName:%s", logName)

	// 查看log文件夹，无则创建
	dir, err := os.ReadDir(logDir)
	if err != nil {
		zap.S().Infof("[InitLogger] %v %v", dir, err)
		return
	}
	logFile.Close() // 关闭之前的log日志

	// 查看是否存在当天日志，无则创建.省略，直接创建新的
	// 创建一个文件输出器,之后不能关闭file
	logFile, err = os.Create(logPath(logName))
	if err != nil {
		panic(err)
	}

	// 配置日志编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 配置控制台输出
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleOutput := zapcore.Lock(os.Stdout)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleOutput, zapcore.InfoLevel)

	// 配置文件输出
	fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileOutput := zapcore.Lock(logFile)
	fileCore := zapcore.NewCore(fileEncoder, fileOutput, zapcore.InfoLevel)

	// 创建一个核心
	core := zapcore.NewTee(consoleCore, fileCore)

	// 创建Logger
	logger := zap.New(core)

	// 将logger设置为全局默认logger
	zap.ReplaceGlobals(logger)
	zap.L().Sync()
}

// 查询log下是否存在当天日志
func isExistLog(path string) bool {
	dir, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	date := time.Now().Format(time.DateTime)
	for i := 0; i < len(dir); i++ {
		if !dir[i].IsDir() &&
			strings.Contains(dir[i].Name(), date) {
			return true
		}
	}
	return false
}

func getLoggerName(srv string) string {
	hostname, _ := os.Hostname()
	date := time.Now().Format(time.DateOnly)
	num := rand.Intn(1000)
	return fmt.Sprintf("%s__%s__%s__%d.log", hostname, date, srv, num)
}

func logPath(name string) string {
	return logDir + "/" + name
}

func InitLogger1() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		zap.S().Warn("[InitLogger]无法启动日志:%s", err.Error())
	}
	zap.ReplaceGlobals(logger)
}
