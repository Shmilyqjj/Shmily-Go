package main

import (
	"fmt"
	"math"
)

func main() {
	// 简短if语句 其中v只在if的作用域
	if v := math.Pow(3, 2); v < 10 {
		fmt.Printf("%v < 10 ", v)
	} else {
		fmt.Printf("%v >= 10 ", v)
	}

}
