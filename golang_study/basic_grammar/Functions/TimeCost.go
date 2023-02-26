package main

import (
	"fmt"
	"time"
)

func main() {
	// 由于 Golang 提供了函数延时执行的功能，借助 defer 我们可以通过函数封装的方式来避免代码冗余
	doSomething(5)
}

func doSomething(times int) {
	defer timeCost(time.Now(), "doSomething")
	for i := 1; i <= times; i++ {
		time.Sleep(time.Duration(1) * time.Second)
	}
}

// timeCost 耗时统计
func timeCost(start time.Time, funcName string) {
	tc := time.Since(start)
	fmt.Printf("Function time cost = %v\n", tc)
}
