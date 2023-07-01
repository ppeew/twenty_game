package initialize

import (
	"encoding/json"
	"user_srv/global"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	global.DEBUG = GetEnvInfo("PPEEW_DEBUG")
	v := viper.New()
	fileName := "config-pro.yaml"
	if global.DEBUG {
		fileName = "config-debug.yaml"
	}
	v.SetConfigFile(fileName)
	err := v.ReadInConfig()
	if err != nil {
		zap.S().Fatalf("[InitConfig]读取nacos配置错误:%s", err.Error())
	}
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		zap.S().Fatalf("[InitConfig]读取nacos配置错误:%s", err.Error())
	}
	//读取yaml文件，去nacos配置中心找配置文件
	clientConfig := constant.ClientConfig{
		NamespaceId: global.NacosConfig.Namespace,
		Username:    global.NacosConfig.User,
		Password:    global.NacosConfig.Password,
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		zap.S().Fatalf("[InitConfig]连接nacos错误:%s", err.Error())
	}
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		zap.S().Fatalf("[InitConfig]读取nacos服务器的配置信息错误:%s", err.Error())
	}
	err = json.Unmarshal([]byte(content), global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("[InitConfig]反序列化ServerConfig失败:%s", err.Error())
	}
}
