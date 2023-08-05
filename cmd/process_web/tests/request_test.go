package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/parnurzeal/gorequest"
)

type Res struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Gender   bool   `json:"gender"`
	Username string `json:"username"`
	Image    string `json:"image"`
}

func TestRequest(t *testing.T) {
	gorequest.New().
		Get("http://139.159.234.134:8000/user/v1/search").
		Set("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NTAsImV4cCI6MTY4OTY5NTQ3OCwiaXNzIjoicHBlZXciLCJuYmYiOjE2ODkyNjM0Nzh9.w4hH23492VGH5aq1b2jVLntFG-gPQnobKthK0lSgSVM").
		Param("id", "54").
		End(func(response gorequest.Response, body string, errs []error) {
			if errs != nil {
				fmt.Println(errs)
			}
			fmt.Println(body)
			fmt.Println("--------")
			fmt.Println(response.Body)
		})
	fmt.Println("ok!")
	time.Sleep(time.Second * 2)
}
