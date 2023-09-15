package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestQps(t *testing.T) {
	// 创建一个等待组来等待所有测试完成
	var wg sync.WaitGroup
	numRequest := 1000
	wg.Add(numRequest)
	start := time.Now()
	for i := 0; i < numRequest; i++ {
		go Request(&wg)
	}

	wg.Wait()

	// 计算测试持续时间
	since := time.Since(start)

	// 计算QPS
	qps := float64(numRequest) / since.Seconds()
	fmt.Printf("QPS: %.2f\n", qps)

}

func Request(wg *sync.WaitGroup) {
	defer wg.Done()
	url := "http://139.159.234.134:8000/game/v1/getRoomList?pageIndex=%3CpageIndex%3E&pageSize=%3CpageSize%3E"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NzU5OSwiZXhwIjoxNjk0Nzk1ODcwLCJpc3MiOiJwcGVldyIsIm5iZiI6MTY5NDM2Mzg3MH0.jImoTOGHto8cTVR-4AXdQX_g2ZD_yppFzRzXlBRndQU")
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
