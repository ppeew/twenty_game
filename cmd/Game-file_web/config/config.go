package config

type JWTConfig struct {
	SigningKey string `json:"key"`
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
	Host        string        `json:"host"`
	Port        int           `json:"port"`
	JWTInfo     JWTConfig     `json:"jwt"`
	ConsulInfo  ConsulConfig  `json:"consul"`
	UserSrvInfo UserSrvConfig `json:"user_srv"`
	MongoInfo   MongoConfig   `json:"mongo"`
	//OssInfo    OssConfig    `json:"ossInfo"`
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
