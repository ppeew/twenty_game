package initialize

import (
	"bufio"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"os"
	"store_web/global"
	"strings"
	"time"
)

// InitTasks 初始化定时任务
func InitTasks() {
	go func() {
		// 隔天创建新日志
		createLogTask := time.NewTicker(time.Hour * 24)
		deleteOldLogTask := time.NewTicker(time.Hour * 24 * 7)
		saveLogTask := time.NewTicker(time.Hour*24 + time.Minute)

		for {
			select {
			case <-createLogTask.C:
				InitLogger()
			case <-deleteOldLogTask.C:
				deleteOldLog()
			case <-saveLogTask.C:
				saveLog()
			}
		}
	}()
}

/*
由admin负责删除所有服务旧日志
每个admin只负责自己主机
请求其他服务删除本地日志接口，类似发信号方式，触发删除功能
做到不干涉服务内部情况
需要获取主机下所有服务，复杂，放弃

改为每个服务自身定时触发删除旧日志
*/
func deleteOldLog() {
	lastWeek := time.Now().AddDate(0, 0, -7).Format(time.DateOnly)

	dir, err := os.ReadDir(logDir)
	if err != nil {
		zap.S().Infof("[deleteOldLog] %v", err)
		return
	}

	for _, file := range dir {
		if !file.IsDir() && isOldLog(file.Name(), lastWeek) {
			err := os.Remove(logPath(file.Name()))
			if err != nil {
				zap.S().Infof("[deleteOldLog] delete is not success:%v", err)
			} else {
				zap.S().Infof("[deleteOldLog] successfully delete %s", file.Name())
			}
		}
	}
}

/*
保存昨天的日志到MongoDB
*/
func saveLog() {
	yesterday := time.Now().AddDate(0, 0, -1).Format(time.DateOnly)
	dir, err := os.ReadDir(logDir)
	if err != nil {
		zap.S().Infof("[deleteOldLog] %v", err)
		return
	}

	for _, file := range dir {
		if !file.IsDir() && getLogDay(file.Name()) == yesterday {
			err := saveInMongo(file.Name())
			if err != nil {
				zap.S().Infof("[saveLog] unsave %s", file.Name())
			} else {
				zap.S().Infof("[saveLog] successfully save %s", file.Name())
			}
		}
	}
}

func saveInMongo(filename string) error {
	logger, err := os.Open(logPath(filename))
	if err != nil {
		return err
	}

	id := 0
	scanner := bufio.NewScanner(logger)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		index := strings.LastIndex(filename, ".")
		global.MongoDB.Collection(filename[:index]).InsertOne(context.Background(), bson.M{
			"_id": id,
			"log": line,
		})
		id++
	}
	return nil
}

func isOldLog(filename, lastWeek string) bool {
	day := getLogDay(filename)
	if day == "" {
		return false
	}
	return day <= lastWeek
}

func getLogDay(name string) string {
	arr := strings.Split(name, "__")
	for _, elem := range arr {
		_, err := time.Parse(time.DateOnly, elem)
		// 成功解析则返回
		if err == nil {
			return elem
		}
	}
	return ""
}
