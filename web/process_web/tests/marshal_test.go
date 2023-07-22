package tests

import (
	"encoding/json"
	"testing"
)

func TestMarshal(t *testing.T) {
	//u:=[]uint{1,2,2,3}
	var user []uint
	user = append(user, 1)
	user = append(user, 2)
	user = append(user, 3)
	c := make([]map[string]interface{}, 0)
	c = append(c, map[string]interface{}{
		"users": user,
	})

	marshal, _ := json.Marshal(c)
	println(string(marshal))
}
