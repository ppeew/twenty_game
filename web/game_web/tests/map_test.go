package tests

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	var m = make(map[int]int)
	//m[3] = 2
	if val, ok := m[3]; ok {
		fmt.Println(val)
	}
}
