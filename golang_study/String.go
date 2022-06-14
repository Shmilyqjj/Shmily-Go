package main

import "fmt"

func main() {
	strConcat()
}

func strConcat() {
	var i int = 123
	var s = "aaa"
	var format_str = "i=%d s=%s"
	var str = fmt.Sprintf(format_str, i, s)
	fmt.Println(str)
}
