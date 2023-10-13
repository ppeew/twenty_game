package initialize

import (
	"time"
)

// InitTasks 初始化定时任务
func InitTasks() {
	go func() {
		// 隔天创建新日志
		createLogTask := time.NewTicker(time.Hour * 24)

		for {
			select {
			case <-createLogTask.C:
				InitLogger()
			}
		}
	}()
}
