package initialize

import (
	"admin_web/global"
	"admin_web/utils"
	"bufio"
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

type server struct {
	name string
	host string
}

var hasNotify = make(map[string]bool)

// InitTasks 初始化定时任务
func InitTasks() {
	go func() {
		// 隔天创建新日志
		createLogTask := time.NewTicker(time.Hour * 24)

		// admin专门任务
		deleteOldLogTask := time.NewTicker(time.Hour * 24 * 7)
		saveLogTask := time.NewTicker(time.Hour*24 + time.Minute)
		monitorServiceTask := time.NewTicker(time.Second * 30)
		monitorHostTask := time.NewTicker(time.Minute * 15)
		resetTask := time.NewTicker(time.Minute * 20)
		for {
			select {
			case <-createLogTask.C:
				InitLogger()
			case <-deleteOldLogTask.C:
				deleteOldLog()
			case <-saveLogTask.C:
				saveLog()
			case <-monitorServiceTask.C:
				monitorService()
			case <-monitorHostTask.C:
				monitorHost()
			case <-resetTask.C:
				resetData()
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

// 监听本机各服务情况
func monitorService() {
	// 创建Consul客户端
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(config)
	if err != nil {
		zap.S().Infof("[monitorService] err:%v", err)
		return
	}

	badService := make([]server, 0)
	for _, checkService := range global.ServerConfig.CheckServices {
		srvs, _, _ := client.Catalog().Service(checkService.Name, "", nil)

		// 检查该服务下所有主机服务是否存在
		for _, host := range checkService.Hosts {
			key := fmt.Sprintf("%s__%s", checkService.Name, host)
			if _, ok := hasNotify[key]; ok {
				continue // 已经通知过的不用再通知
			}

			has := false
			for _, srv := range srvs {
				if host == srv.ServiceAddress {
					has = true
					break
				}
			}

			if !has {
				badService = append(badService, server{
					name: checkService.Name,
					host: host,
				})
			}
		}
	}

	if len(badService) > 0 {
		notifyMailBox(badService)
	}
}

// 监听主机情况
func monitorHost() {
	if _, ok := hasNotify[global.ServerConfig.Host]; ok {
		return // 已经通知过的不用再通知
	}

	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		zap.S().Errorf("Failed to get CPU usage: %v", err)
		return
	}

	memUsage, err := mem.VirtualMemory()
	if err != nil {
		zap.S().Errorf("Failed to get memory usage: %v", err)
		return
	}

	zap.S().Infof("CPU Usage: %.2f% %", cpuUsage[0])
	zap.S().Infof("Memory Usage: %.2f% %", memUsage.UsedPercent)

	msg := ""
	if cpuUsage[0] > 80 {
		msg += "主机" + global.ServerConfig.Host + " CPU使用率超过80% </br>"
	}

	if memUsage.UsedPercent > 80 {
		msg += "主机" + global.ServerConfig.Host + " 内存使用率超过80% </br>"
	}

	err = utils.SendMessage(global.ServerConfig.SendMailBox, msg)
	if err != nil {
		zap.S().Infof("[notifyMailBox] 发送邮箱信息失败 %s", err.Error())
	} else {
		hasNotify[global.ServerConfig.Host] = true
	}
}

func notifyMailBox(services []server) {
	msg := ""
	for _, srv := range services {
		fmt.Printf("%+v\n", srv)
		msg += fmt.Sprintf("服务名：%s </br> 服务地址：%s </br>", srv.name, srv.host)
	}
	msg += "以上服务已在Consul中消失，请排查原因</br>"
	err := utils.SendMessage(global.ServerConfig.SendMailBox, msg)
	if err != nil {
		zap.S().Infof("[notifyMailBox] 发送邮箱信息失败 %s", err.Error())
		return
	} else {
		for _, srv := range services {
			key := fmt.Sprintf("%s__%s", srv.name, srv.host)
			hasNotify[key] = true
		}
	}
}

func resetData() {
	hasNotify = make(map[string]bool)
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
