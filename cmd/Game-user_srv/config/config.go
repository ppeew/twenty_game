package config

type MysqlConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      uint64 `mapstructure:"port" json:"port"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
	DataId    string `mapstructure:"dataid" json:"dataid"`
	Group     string `mapstructure:"group" json:"group"`
}

type ConsulConfig struct {
	Name string `json:"name"` //服务在注册中心的名字
	Host string `json:"host"`
	Port int    `json:"port"`
}

type GameSrvConfig struct {
	Name string `json:"name"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

type ServerConfig struct {
	Host        string        `json:"host"`
	MysqlInfo   MysqlConfig   `json:"mysql"`
	ConsulInfo  ConsulConfig  `json:"consul"`
	GameSrvInfo GameSrvConfig `json:"game_srv"`
	RedisInfo   RedisConfig   `json:"redis"`
	MongoInfo   MongoConfig   `json:"mongo"`
}

type MongoConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}
