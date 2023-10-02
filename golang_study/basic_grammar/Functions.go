package main

import "fmt"

func main() {
	a := 1
	b := 2
	println(opt("+")(a, b))
	println(opt("-")(a, b))

	// Function接收一个指针
	d := &Data{
		Id: 1,
	}
	d.Init()
	println(d.C)
}

// 函数也是一个值 可以用作函数的返回值 可以像值一样传递
func opt(optStr string) func(x int, y int) int {
	switch optStr {
	case "+":
		return func(x int, y int) int {
			return x + y
		}
	case "-":
		return func(x int, y int) int {
			return x - y
		}
	default:
		return func(x int, y int) int {
			return x + y
		}
	}
}

type Data struct {
	Id   int
	Name string
	C    string
}

// 函数接收一个指针: 可以避免在每次调用方法时复制该值 如果接收一个结构体则更高效,可以联想为Java的类变量
func (d *Data) Init() {
	if d.Name == "" {
		d.Name = "abc"
	}
	println(d.Id)
	println(d.Name)
	d.C = fmt.Sprintf("%d=%s", d.Id, d.Name)
}
