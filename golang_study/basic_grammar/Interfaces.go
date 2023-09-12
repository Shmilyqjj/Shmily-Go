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

func main() {
	rect := Rectangle{Width: 5, Height: 3}
	circle := Circle{Radius: 2}

	// 调用方法，可以传递不同类型的参数
	CalculateArea(rect)   // 输出: Area: 15.00
	CalculateArea(circle) // 输出: Area: 12.56
}
