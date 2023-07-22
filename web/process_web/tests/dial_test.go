package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestDialWS(t *testing.T) {
	token := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NCwiZXhwIjoxNjgxOTc0OTUwLCJpc3MiOiJwcGVldyIsIm5iZiI6MTY4MTg4ODU1MH0.VusRj-1ZTfDclY5OHPOqGbypPiMeGT4S1sYcCgI9uqU",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NSwiZXhwIjoxNjgxOTc1NDY2LCJpc3MiOiJwcGVldyIsIm5iZiI6MTY4MTg4OTA2Nn0.DufQD563sxzXGygNUe-Krm0mTWYI9AVOIHpkYFZqEbQ",
	}
	for i, t := range token {
		go func() {
			dial, _, err := websocket.DefaultDialer.Dial("ws://localhost:8083/v1/userIntoRoom?room_id=1234&/ws", http.Header{
				"token": []string{t},
			})
			if err != nil {
				panic(err)
			}
			for {
				_, p, err := dial.ReadMessage()
				if err != nil {
					panic(err)
				}
				fmt.Printf("用户%d收到:%s\n", i, string(p))
			}
			dial.Close()
		}()
		time.Sleep(1000000 * time.Second)
	}
}
