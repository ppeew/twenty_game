package config

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type ConsulConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Name string `json:"name"`
}

type UserSrvConfig struct {
	Name string `json:"name"`
}

type ServerConfig struct {
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	JWTInfo     JWTConfig     `mapstructure:"jwt"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv"`
	MongoInfo   MongoConfig   `mapstructure:"mongo"`
	//OssInfo    OssConfig    `mapstructure:"ossInfo"`
}

type OssConfig struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
	Host      string `json:"host"`
	UploadDir string `json:"uploadDir"`
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

type MongoConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}
