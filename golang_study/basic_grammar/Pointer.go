package main

// Golang指针
// 在Go中 * 代表取指针地址中存的值
// 在Go中 & 代表取一个值的地址 (取地址运算符&，当放置在变量前面时，可以获取该变量的内存地址)
// 指针储存的是一个值的地址，但本身这个指针也需要地址来储存

import "fmt"

func main() {
	var p *int      //定义一个指针p
	p = new(int)    // 创建一块内存分配给p
	*p = 1          //
	fmt.Println(p)  // p是一个指针，指针的值为内存地址p
	fmt.Println(*p) // p这个内存地址中存的值为1
	fmt.Println(&p) // 指针p的内存地址
	errorExample()
	fmt.Println("===========================")
	x := 10
	i := changeX(&x) // x是值 入参是&x即取x的地址
	fmt.Println(*i)  // 输出时取地址中存的值
	i1 := changeX(i) // 地址可以直接做入参传入
	fmt.Println(*i1)
	fmt.Println(*changeX(new(int))) // new(int) 开辟一块地址 存的值为0
	fmt.Println(*XAdd(&x, 1))
}

func errorExample() {
	// 初始化指针的时候会将指针i的值赋为nil
	// 但 i 的值代表的是 *i 的地址， nil 的话系统还并没有给 *i 分配地址，所以这时给 *i 赋值肯定会出错
	var i *int
	i = new(int) // 在给指针赋值前可以先创建一块内存分配给赋值对象即可
	*i = 1
	fmt.Println(i)
	fmt.Println(*i)
	fmt.Println(&i)
}

func changeX(x *int) *int {
	// 通过使用指针，我们可以传递变量的引用（例如，作为函数的参数），而不是传递变量的副本，从而减少内存使用量并提高效率。
	*x += 6  // 传入的参数是地址 *x取到值并加6
	return x // 返回x的地址
}

func XAdd(x *int, y int) *int {
	*x = *x + y
	return x
}
