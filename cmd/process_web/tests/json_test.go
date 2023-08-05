package tests

import (
	"encoding/json"
	"testing"
	"time"
)

type User struct {
	Ready bool `json:"Ready"`
}

type UsersInfos struct {
	ID    uint32  `json:"Id"`
	Users []*User `json:"Users"`
}

func TestJsonMarshal(t *testing.T) {
	var users []*User
	users = append(users, &User{Ready: false})
	users = append(users, &User{Ready: true})
	infos := UsersInfos{
		ID:    1,
		Users: users,
	}
	marshal, _ := json.Marshal(infos)
	println(string(marshal))
}

func TestJsonUnmarshal(t *testing.T) {
	var infos UsersInfos
	str := `{"Id":1,"Users":[{"Ready":false},{"Ready":true}]}`
	//注意三种区别:
	//1. ``
	//2. ""
	//3 ''
	_ = json.Unmarshal([]byte(str), &infos)
	time.Sleep(12 * time.Second)
}

type People struct {
	Name string `json:"Name"`
	IsOk bool   `json:"IsOk"`
}

type All struct {
	Peoples []*People
	Id      uint32
}

func TestJsonUnmarshalInterface(t *testing.T) {
	//错误用法
	//data := People{Name: "hhhhh"}
	//bytes, _ := json.Marshal(data)
	//aa := map[string]interface{}{
	//	"data": string(bytes),
	//}
	//marshal, _ := json.Marshal(aa)
	//println(string(marshal)) // {"data":"{\"Name\":\"hhhhh\"}"}
	//可见，不要多次序列化，否则会存在转义字符

	//正常的用法,直接赋值接口没序列化的内容
	//data := People{Name: "hh2h2hh", IsOk: true}
	//var datas []*People
	//datas = append(datas, &data)
	//
	//fmt.Printf("%+v\n", room)
	//aa := map[string]interface{}{
	//	"data": room,
	//}
	//marshal, _ := json.Marshal(aa)
	//println(string(marshal)) // {"data":{"Name":"hhhhh"}}
}
