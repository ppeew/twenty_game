package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"sync"
	"testing"
)

func TestSystem(t *testing.T) {
	//测试系统并发情况，支持10000并发
	group := sync.WaitGroup{}
	group.Add(10000)
	for i := 0; i < 10000; i++ {
		go func(i int) {
			defer group.Done()
			//1.新建用户
			token, err := RegistUser(i)
			if err != nil {
				fmt.Println(fmt.Sprintf("%d线程注册失败", i))
				return
			}
			//2.创建房间
			err = CreateRoom(token, i)
			if err != nil {
				fmt.Println(fmt.Sprintf("%d线程创房失败", i))
				return
			}
			//3.开始游戏
			//BeginGame()
			//4.游戏结束
		}(i)
	}
	group.Wait()
}

func Test1(t *testing.T) {
	RegistUser(1)
}

func RegistUser(i int) (string, error) {

	url := "http://8.134.163.22:8000/user/v1/register"
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("nickname", "test!")
	_ = writer.WriteField("gender", "true")
	_ = writer.WriteField("username", fmt.Sprintf("%d", 100+i))
	_ = writer.WriteField("password", "1234567")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	m := make(map[string]interface{})
	json.Unmarshal(body, &m)
	fmt.Println(m["token"])
	return m["token"].(string), nil
}

func CreateRoom(token string, i int) error {
	url := "http://8.134.163.22:8000/game/v1/createRoom"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("room_id", fmt.Sprintf("%d", i))
	_ = writer.WriteField("max_user_number", "1")
	_ = writer.WriteField("game_count", "5")
	_ = writer.WriteField("room_name", "测试房间")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("token", token)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
