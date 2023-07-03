package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Warnf("初始化[InitSentinel]异常:%s", err)
		return
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "user_web",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              100,
			//WarmUpPeriodSec:        60, // 60s内逐渐达到最大1000并发
			StatIntervalInMs: 1000,
		},
	})
	if err != nil {
		zap.S().Warnf("[InitSentinel]:%s", err)
		return
	}
}
