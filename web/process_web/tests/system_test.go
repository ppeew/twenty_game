package tests

import (
	"bytes"
	"encoding/json"
	"errors"
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
	"sync/atomic"
	"testing"
	"time"
)

var isLocal = false

func Test1(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTU2LCJleHAiOjE2OTA4MDM0NjUsImlzcyI6InBwZWV3IiwibmJmIjoxNjkwMzcxNDY1fQ.hZt6hWO7Hb-7NvoSKInKA1qy_OOvC5qv2r0M6QHGmHY"
	roomID := "11"
	EnterRoom(token, roomID)
}

func TestSystem(t *testing.T) {
	//测试系统并发情况，支持10000并发
	group := sync.WaitGroup{}
	testNum := 1000
	group.Add(testNum)

	var count int64 = 0
	go countNum(&count)
	for i := 0; i < testNum; i++ {
		go func(i int) {
			//1.新建用户
			userName := randString(5)
			token, err := RegisterUser(userName)
			if err != nil {
				// 注册失败则尝试登录获取token
				token, err = Login(userName)
				if err != nil {
					return
				}
			}

			//2.创建房间
			err = CreateRoom(token, userName)
			if err != nil {
				fmt.Println(fmt.Sprintf("%d线程创房失败", i))
				return
			}

			atomic.AddInt64(&count, 1)
			//3.进入房间
			time.Sleep(time.Second)
			wsConn, err := EnterRoom(token, userName)
			if err != nil {
				fmt.Println(fmt.Sprintf("%d线程进入房间失败", i))
				return
			}

			//4.开始游戏
			BeginGame(wsConn)

			//5.游戏结束
			EndGame(wsConn, userName)

			//6.退出房间
			QuitRoom(wsConn, userName)
			group.Done()
		}(i)
	}
	group.Wait()
}

func countNum(counter *int64) {
	// 启动定时器
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	// 定时器的回调函数
	go func() {
		for range ticker.C {
			// 输出当前并发数
			fmt.Printf("[当前并发数]: %d\n", atomic.LoadInt64(counter))
			// 重置计数器
			atomic.StoreInt64(counter, 0)
		}
	}()

	time.Sleep(time.Minute * 3)
}

func randString(len int) string {
	num := 1
	for i := 0; i < len-1; i++ {
		num *= 10
	}
	//num + rand.Intn(num*10-num)
	return strconv.Itoa(rand.Intn(num))
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

func RegisterUser(i string) (string, error) {
	url := Host() + "/user/v1/register"
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("nickname", "test!")
	_ = writer.WriteField("gender", "true")
	_ = writer.WriteField("username", i)
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
	if res.StatusCode != http.StatusOK {
		return "", errors.New("注册失败")
	}

	m, _ := getBody(res.Body)
	fmt.Printf("[Register] %s注册成功\n", i)
	return m["token"].(string), nil
}

func Login(userName string) (string, error) {
	url := Host() + "/user/v1/login"
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("username", userName)
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
	if res.StatusCode != http.StatusOK {
		return "", errors.New("登录失败")
	}

	m, _ := getBody(res.Body)
	fmt.Println("[Login] 登录成功")
	return m["token"].(string), nil
}

func CreateRoom(token string, i string) error {
	url := Host() + "/game/v1/createRoom"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("room_id", i)
	_ = writer.WriteField("max_user_number", "1")
	_ = writer.WriteField("game_count", "1")
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
	fmt.Printf("[CreateRoom] %+v\n", m)
	return nil
}

func getTargetAddress(token string, roomID string) (string, error) {
	url := Host() + "/game/v1/selectRoomServer?room_id=" + roomID
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
	if res.StatusCode != http.StatusOK {
		return "", errors.New("获取远程服务器失败")
	}

	m, _ := getBody(res.Body)
	fmt.Printf("[EnterRoom] %+v\n", m)
	return m["server"].(map[string]interface{})["serverInfo"].(string), nil
}

func connectSocket(targetURL string, token string, roomID string) *websocket.Conn {
	url := fmt.Sprintf("ws://%s/v1/connectSocket?token=%s&room_id=%s", targetURL, token, roomID)
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// 创建 WebSocket 连接
	wsDialer := websocket.DefaultDialer
	wsConn, _, err := wsDialer.Dial(req.URL.String(), nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// 在 WebSocket 连接上发送消息
	message := healthMsg()
	fmt.Printf("Send message: %s\n", message)
	err = wsConn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// 接收 WebSocket 服务器端发送的消息
	_, message, err = wsConn.ReadMessage()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Printf("Received message: %s\n", message)
	return wsConn
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
		"type": 204,
	})
	return res
}

func quitRoomMsg() []byte {
	res, _ := json.Marshal(map[string]interface{}{
		"type": 200,
	})
	return res
}

func EnterRoom(token string, roomID string) (*websocket.Conn, error) {
	targetURL, err := getTargetAddress(token, roomID)
	wsConn := connectSocket(targetURL, token, roomID)
	return wsConn, err
}

func BeginGame(wsConn *websocket.Conn) {
	msg := beginGameMsg()
	fmt.Printf("[BeginGame] %s\n", msg)
	wsConn.WriteMessage(websocket.TextMessage, msg)
}

func EndGame(wsConn *websocket.Conn, userName string) {
	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			break
		}
		m, err := getBody(bytes.NewReader(message))
		fmt.Printf("Received message: %f\n", m["msgType"])
		if int(m["msgType"].(float64)) == 304 {
			fmt.Printf("%s 游戏结束\n", userName)
			break
		}
	}
}

func QuitRoom(wsConn *websocket.Conn, userName string) {
	msg := quitRoomMsg()
	fmt.Printf("[QuitRoom] %s退出房间\n", userName)
	wsConn.WriteMessage(websocket.TextMessage, msg)
}
