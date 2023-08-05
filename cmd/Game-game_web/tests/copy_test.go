package tests

import (
	"testing"
	"time"
)

type a struct {
	m map[int]int
}

type b struct {
	val a
}

func TestCopy(t *testing.T) {
	go func() {
		c := make(map[int]int)
		t1 := a{
			m: c,
		}
		c[2] = 1
		go func(t1 a) {
			time.Sleep(2 * time.Second)
			t2 := b{
				val: t1,
			}
			if _, exist := t2.val.m[2]; exist {
				println("我已经有了")
			}
		}(t1)
	}()
	time.Sleep(10 * time.Second)
}
