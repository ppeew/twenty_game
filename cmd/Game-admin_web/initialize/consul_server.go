package initialize

import (
	"admin_web/global"
	"github.com/hashicorp/consul/api"
	"time"
)

func GetConsulServer() {
	tick := time.Tick(time.Second * 10)
	config := api.DefaultConfig()
	config.Address = "139.159.234.134:8500"
	c, _ := api.NewClient(config)
	for {
		select {
		case <-tick:
			services, _, err := c.Catalog().Service("admin-web", "", nil)
			if err == nil {
				global.ConsulHallWebServices = services
			}
		}
	}
}
