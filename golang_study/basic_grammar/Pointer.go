package main

// Golang指针
// 在 Go 中 * 代表取指针地址中存的值，& 代表取一个值的地址
// 指针储存的是一个值的地址，但本身这个指针也需要地址来储存

import "fmt"

func main() {
	var p *int      //定义一个指针p
	p = new(int)    // 创建一块内存分配给p
	*p = 1          //
	fmt.Println(p)  // p是一个指针，指针的值为内存地址p
	fmt.Println(*p) // p这个内存地址中存的值为1
	fmt.Println(&p) // p这个指针的内存地址

	errorExample()
}

func errorExample() {
	// go 初始化指针的时候会为指针 i 的值赋为 nil ，但 i 的值代表的是 *i 的地址， nil 的话系统还并没有给 *i 分配地址，所以这时给 *i 赋值肯定会出错

	var i *int
	//i = new(int)   // 解决：在给指针赋值前可以先创建一块内存分配给赋值对象即可
	*i = 1
	fmt.Println(i)
	fmt.Println(*i)
	fmt.Println(&i)
}
