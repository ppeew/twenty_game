package main

import (
	"fmt"
	"time"
)

func mm1() {
	for {
		for i := 0; i < 10; i++ {
			ch8 <- 8 * i
		}
		time.Sleep(2 * time.Second)
	}
	//close(ch8)
}
func mm2() {
	//time.Sleep(2 * time.Second)
	for data := range ch8 {
		fmt.Print(data, "\t")
	}
	fmt.Printf("读取完成\n")
}

var ch8 = make(chan int, 1024)

func main() {
	go mm1()
	mm2()
	mm2()
	time.Sleep(10 * time.Second)

}
