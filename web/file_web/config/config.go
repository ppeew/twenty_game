package config

type JWTConfig struct {
	SigningKey string `json:"key"`
}

type ConsulConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type ServerConfig struct {
	Port       int          `json:"port"`
	JWTInfo    JWTConfig    `json:"jwt"`
	ConsulInfo ConsulConfig `json:"consul"`
	OssInfo    OssConfig    `json:"ossInfo"`
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
