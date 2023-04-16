package config

type GameSrvConfig struct {
	Name string `json:"name"`
}

type JWTConfig struct {
	SigningKey string `json:"key"`
}

type ConsulConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

type ServerConfig struct {
	Port        int           `json:"port"`
	GameSrvInfo GameSrvConfig `json:"game_srv"`
	JWTInfo     JWTConfig     `json:"jwt"`
	ConsulInfo  ConsulConfig  `json:"consul"`
	RedisInfo   RedisConfig   `json:"redis"`
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
