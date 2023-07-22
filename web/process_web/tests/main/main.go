package main

import (
	"encoding/json"
	"fmt"
)

type Student struct {
	Name  string
	Age   int
	Hobby string
}

// struct序列化
func structJson() {
	student := Student{
		Name:  "zhangsan",
		Age:   20,
		Hobby: "coding",
	}
	data, err := json.Marshal(&student)
	if err != nil {
		fmt.Printf("err = %v", err)
	}
	fmt.Printf("struct serialized ：%v\n", string(data))
}

// map序列化
func mapJson() {
	var a map[string]interface{}
	a = make(map[string]interface{})
	a["name"] = "lisi"
	a["age"] = 22
	a["hobby"] = "reading"
	data, err := json.Marshal(a)
	if err != nil {
		fmt.Printf("err = %v", err)
	}
	fmt.Printf("map serialized ：%v\n", string(data))
}

// slice序列化
func sliceJson() {
	var slice []map[string]interface{}
	var m1 map[string]interface{}
	m1 = make(map[string]interface{})
	m1["name"] = "wangwu"
	m1["age"] = 26
	m1["hobby"] = "play games"
	slice = append(slice, m1)

	var m2 map[string]interface{}
	m2 = make(map[string]interface{})
	m2["name"] = "zhaoliu"
	m2["age"] = 28
	m2["hobby"] = "watching tv"
	slice = append(slice, m2)
	data, err := json.Marshal(slice)
	if err != nil {
		fmt.Printf("err = %v", err)
	}
	fmt.Printf("slice serialized ：%v\n", string(data))

}
func main() {
	structJson()
	//mapJson()
	//sliceJson()
}
