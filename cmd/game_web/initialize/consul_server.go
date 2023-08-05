package initialize

import (
	"game_web/global"
	"github.com/hashicorp/consul/api"
	"time"
)

func GetConsulServer() {
	tick := time.Tick(time.Second * 10)
	config := api.DefaultConfig()
	config.Address = "139.159.234.134:8500"
	c, _ := api.NewClient(config)
	select {
	case <-tick:
		services, _, err := c.Catalog().Service("process-web", "", nil)
		if err == nil {
			global.ConsulProcessWebServices = services
		}
	}
}
