package main

import (
	"fmt"
)

//使用Interface使方法兼容不同数据类型

// 定义一个接口
type Shape interface {
	Area() float64
}

// 定义一个矩形类型
type Rectangle struct {
	Width  float64
	Height float64
}

// 矩形类型实现 Shape 接口的 Area 方法
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// 定义一个圆形类型
type Circle struct {
	Radius float64
}

// 圆形类型实现 Shape 接口的 Area 方法
func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

// 使用 Shape 接口作为参数的方法
func CalculateArea(s Shape) {
	area := s.Area()
	fmt.Printf("Area: %.2f\n", area)
}

// ===============================

// 接口也是值。它们可以像其它值一样传递。 接口值可以用作函数的参数或返回值。
// 在内部，接口值可以看做包含值和具体类型的元组： (value, type)
// 接口值调用方法时会执行其底层类型的同名方法。
type Value interface {
	PrintValue()
}
type S struct {
	S string
}
type I struct {
	I int
}

func (s S) PrintValue() {
	println(s.S)
}
func (i I) PrintValue() {
	println(i.I)
}
func describe(v interface{}) {
	fmt.Printf("(%v %T) \n", v, v)
}

func main() {
	rect := Rectangle{Width: 5, Height: 3}
	circle := Circle{Radius: 2}
	// 调用方法，可以传递不同类型的参数
	CalculateArea(rect)   // 输出: Area: 15.00
	CalculateArea(circle) // 输出: Area: 12.56

	// ==============
	i := S{"aaa"}
	describe(i)
	j := I{111}
	describe(j)
	describe(nil)
	i.PrintValue()
	j.PrintValue()

	// ======空接口可以保存任何类型的值========
	var x interface{}
	x = 1
	x = I{11}
	x = "aa"
	describe(x)

	// 类型断言 提供了访问接口值底层具体值的方式。 t, ok := i.(T)  或  t := i.(T) 后者可能触发panic
	r := x.(string)
	fmt.Println(r)
	r, ok := x.(string)
	fmt.Println(r, ok)
	//l := x.(int)
	//fmt.Println(l)
	l, ok := x.(int)
	fmt.Println(l, ok)

	// Interface值类型选择
	switch x.(type) {
	case string:
		println("string")
	case int:
		println("int")
	case bool:
		println("boolean")
	default:
		println("other type")
	}
}
