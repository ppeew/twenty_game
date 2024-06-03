package tests

import (
	"fmt"
	"testing"
)

type Node struct {
	val  int
	pre  *Node
	next *Node
}

type List struct {
	head *Node
	len  int
}

func MakeList(i int) *List {
	len := i
	head := &Node{
		val:  0,
		pre:  nil,
		next: nil,
	}
	cur := head
	for i != 0 {
		node := &Node{
			val:  i,
			pre:  head,
			next: nil,
		}
		(*cur).next = node
		cur = node
		i--
	}
	return &List{
		head: head,
		len:  len,
	}
}

func Sort(list *List) *Node {
	// 0 10 9 8 7 6 5 4 3 2 1
	size := list.len
	head := list.head

	sortNum := 0
	for i := 1; i < size+1; i++ {
		// i指定节点位置
		for j := i + 1; j < size-sortNum+1; j++ {
			// 要找到i位置的node和j位置node
			pf1 := findNode(list, i-1)
			f1 := pf1.next
			f2 := f1.next
			// swap
			if f1.val > f2.val {
				t := f2.next
				head.next = f2
				f2.next = f1
				f1.pre = f2
				f1.next = t
			}
		}
		sortNum++
	}
	return head
}

func findNode(list *List, n int) *Node {
	if n == 0 {
		return list.head
	}
	slow := list.head
	fast := list.head
	for i := 0; i < n; i++ {
		fast = fast.next
	}

	for fast.next != nil {
		slow = slow.next
		fast = fast.next
	}
	return slow
}

func TestDecode(t *testing.T) {
	list := MakeList(10)
	Sort(list)
	for list.head != nil {
		fmt.Println(list.head)
		list.head = list.head.next
	}
}
