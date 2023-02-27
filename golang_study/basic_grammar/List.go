package main

import (
	"container/list"
	"fmt"
)

func main() {
	listCopy()
}

// listCopy list的
func listCopy() {
	l := list.New()
	l.PushBack("1")
	l.PushBack("2")
	l.PushBack("3")
	l.PushBack("4")
	l.PushBack("5")
	l.PushBack("6")
	l.PushBack("7")
	l.PushBack("8")

	// 拷贝地址
	lCopyAddr := l
	// deep copy
	lCopyData := list.New()
	lCopyData.PushBackList(l)

	// 初始化 清空列表
	l.Init()

	fmt.Printf("l: %v \nlCopyAddr: %v \nlCopyData: %v \n", l, lCopyAddr, lCopyData)
	fmt.Printf("l: %d \nlCopyAddr: %d \nlCopyData: %d \n", l.Len(), lCopyAddr.Len(), lCopyData.Len())

}
