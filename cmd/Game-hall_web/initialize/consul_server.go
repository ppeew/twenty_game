package initialize

import (
	"github.com/hashicorp/consul/api"
	"hall_web/global"
	"time"
)

func GetConsulServer() {
	tick := time.Tick(time.Second * 3)
	config := api.DefaultConfig()
	config.Address = "139.159.234.134:8500"
	c, _ := api.NewClient(config)
	for {
		select {
		case <-tick:
			services, _, err := c.Catalog().Service("hall-web-dev", "", nil)
			if err == nil {
				global.ConsulHallWebServices = services
			}
		}
	}
}
