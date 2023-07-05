package tests

import (
	"fmt"
	"sync"
	"testing"
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
	time.Sleep(2 * time.Second)
	for data := range ch8 {
		fmt.Print(data, "\t")
	}
	fmt.Printf("读取完成\n")
}

var ch8 = make(chan int, 1024)

func TestChan(t *testing.T) {
	go mm1()
	mm2()
	mm2()
	time.Sleep(10 * time.Second)

}

func TestChan1(t *testing.T) {
	c := make(chan int)
	//var c chan int
	wg := &sync.WaitGroup{}
	wg.Add(2)
	count := 10000
	go func() {
		for i := 0; i < count; i++ {
			c <- i
			println("write ", i)
		}
		println("write ending")
		wg.Done()
	}()

	time.Sleep(time.Second * 3)
	go func() {
		for {
			num := <-c
			if num == count-1 {
				wg.Done()
			}
			println("num: ", num)
			time.Sleep(time.Second)
		}
	}()
	wg.Wait()
}

func TestChan2(t *testing.T) {
	//var c map[uint32]chan bool
	var c = make(chan bool)
	//go func() {
	//	for true {
	//		time.Sleep(time.Second)
	//		//fmt.Print(<-c)
	//	}
	//}()

	go func() {
		for true {
			//time.Sleep(time.Second)
			c <- true
			//fmt.Print("111")
		}
	}()
	time.Sleep(time.Second * 5)
}
