package config

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type ConsulConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Name string `json:"name"`
}

type ServerConfig struct {
	Host          string         `mapstructure:"host"`
	Port          int            `mapstructure:"port"`
	JWTInfo       JWTConfig      `mapstructure:"jwt"`
	ConsulInfo    ConsulConfig   `mapstructure:"consul"`
	MysqlInfo     MysqlConfig    `mapstructure:"mysql"`
	RedisInfo     RedisConfig    `mapstructure:"redis"`
	MongoInfo     MongoConfig    `mapstructure:"mongo"`
	CheckServices []CheckService `mapstructure:"checkServices"`
	SendMailBox   []string       `mapstructure:"sendMailBox"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

type MysqlConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type MongoConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type CheckService struct {
	Name  string   `json:"name"`
	Hosts []string `json:"hosts"`
}
