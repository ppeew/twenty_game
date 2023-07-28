package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"
)

var isLocal = false
var wsConn = new(websocket.Conn)

func Test1(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTU2LCJleHAiOjE2OTA4MDM0NjUsImlzcyI6InBwZWV3IiwibmJmIjoxNjkwMzcxNDY1fQ.hZt6hWO7Hb-7NvoSKInKA1qy_OOvC5qv2r0M6QHGmHY"
	roomID := 11
	EnterRoom(token, roomID)
}

func TestSystem(t *testing.T) {
	//测试系统并发情况，支持10000并发
	group := sync.WaitGroup{}
	testNum := 10
	for i := 0; i < testNum; i++ {
		group.Add(testNum)
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
			//3.进入房间
			err = EnterRoom(token, i)
			if err != nil {
				fmt.Println(fmt.Sprintf("%d线程进入房间失败", i))
				return
			}

			//4.开始游戏
			BeginGame()

			//5.游戏结束
			group.Done()
		}(i)
	}
	group.Wait()
}

func randString() string {
	i := strconv.Itoa(rand.Intn(10000))
	return time.Now().String() + i
}

func getBody(r io.Reader) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	m := make(map[string]interface{})
	json.Unmarshal(body, &m)
	return m, nil
}

func Host() string {
	if isLocal {
		return "http://127.0.0.1:8000"
	}
	return "http://139.159.234.134:8000"
}

func RegistUser(i int) (string, error) {
	url := Host() + "/user/v1/register"
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("nickname", "test!")
	_ = writer.WriteField("gender", "true")
	_ = writer.WriteField("username", randString())
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

	m, _ := getBody(res.Body)
	//fmt.Println(m["token"])
	return m["token"].(string), nil
}

func CreateRoom(token string, i int) error {
	url := Host() + "/game/v1/createRoom"
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

	m, _ := getBody(res.Body)
	fmt.Printf("%+v\n", m)
	return nil
}

func getTargetAddress(token string, roomID int) (string, error) {
	url := Host() + "/game/v1/selectRoomServer?room_id=" + fmt.Sprintf("%d", roomID)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("token", token)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	m, _ := getBody(res.Body)
	fmt.Printf("%+v\n", m)
	return m["server"].(map[string]interface{})["serverInfo"].(string), nil
}

func connectSocket(targetURL string, token string, roomID int) {
	url := fmt.Sprintf("ws://%s/v1/connectSocket?token=%s&room_id=%d", targetURL, token, roomID)
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建 WebSocket 连接
	wsDialer := websocket.DefaultDialer
	wsConn, _, err = wsDialer.Dial(req.URL.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// 在 WebSocket 连接上发送消息
	message := healthMsg()
	log.Printf("Send message: %s", message)
	err = wsConn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Fatal(err)
	}

	// 接收 WebSocket 服务器端发送的消息
	_, message, err = wsConn.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Received message: %s", message)
}

func healthMsg() []byte {
	d := map[string]interface{}{
		"userID": 2,
		"data":   "ok",
	}

	res, _ := json.Marshal(map[string]interface{}{
		"type":        101,
		"chatMsgData": d,
	})
	return res
}

func beginGameMsg() []byte {
	res, _ := json.Marshal(map[string]interface{}{
		"type": 24,
	})
	return res
}

func EnterRoom(token string, roomID int) error {
	targetURL, err := getTargetAddress(token, roomID)
	connectSocket(targetURL, token, roomID)
	return err
}

func BeginGame() {
	wsConn.WriteMessage(websocket.TextMessage, beginGameMsg())
}

func EndGame() {
	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			break
		}
		log.Printf("Received message: %s", message)
		m, err := getBody(bytes.NewReader(message))
		if m["type"]
	}
}
