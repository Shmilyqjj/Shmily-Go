package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	// go的switch与java不同 匹配到第一个即截止 (默认break)
	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "linux":
		fmt.Println("Linux.")
	case "darwin":
		fmt.Println("OS X.")
	default:
		fmt.Printf("%s.\n", os)
	}

	// 没条件的switch
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
}
