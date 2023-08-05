package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/consul/api"
)

func TestConsulFind(t *testing.T) {
	config := api.DefaultConfig()
	config.Address = "139.159.234.134:8500"
	c, _ := api.NewClient(config)
	service, _, err := c.Catalog().Service("user-srv", "", nil)
	if err != nil {
		panic(err)
	}
	for _, s := range service {
		fmt.Printf("%s,%d\n", s.ServiceAddress, s.ServicePort)
	}
}
