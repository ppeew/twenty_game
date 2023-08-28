package tests

import (
	"sync"
	"testing"
)

type ccc struct {
	ch chan struct{}
}

func TestSyncMap(t *testing.T) {
	var m sync.Map
	value, ok := m.Load("123")
	println(ok)
	ws := value.(*ccc)
	<-ws.ch
}
