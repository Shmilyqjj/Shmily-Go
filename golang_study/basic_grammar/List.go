package main

import (
	"container/list"
	"fmt"
	"strings"
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

	// 数组初始化固定值
	strs := []string{"a", "b", "c", "d", "e", "f"}
	println(strings.Join(strs, ","))

	// 数组切片
	var strs1 = strs[1:4] // 下标1-3的所有元素
	println(strings.Join(strs1, ","))

	// 切片就像数组的引用 切片并不存储任何数据，它只是描述了底层数组中的一段。
	//更改切片的元素会修改其底层数组中对应的元素。
	//与它共享底层数组的切片都会观测到这些修改。
	strs1[0] = "x"
	strs1[1] = "y"
	println(strings.Join(strs1, ","))

	// 切片下界的默认值为 0，上界则是该切片的长度。
	var strs2 = strs[:3]
	var strs3 = strs[1:]
	fmt.Println(strs2)
	fmt.Println(strs3)
	fmt.Printf("len=%d, cap=%d \n", len(strs), cap(strs))

	// 使用make创建切片
	strs4 := make([]string, 4)
	fmt.Printf("len=%d, cap=%d \n", len(strs4), cap(strs4))

	// 向数组/切片追加元素
	strs4 = append(strs4, "a")
	fmt.Println(strs4)
}
