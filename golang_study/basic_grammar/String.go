package main

import (
	"fmt"
	"strings"
)

func main() {
	strConcat()
	splitStr()
	replaceStr()
}

func strConcat() {
	var i int = 123
	var s = "aaa"
	var format_str = "i=%d s=%s"
	var str = fmt.Sprintf(format_str, i, s)
	fmt.Println(str)
}

func splitStr() {
	s := "ba:bbb:aaa"
	if strings.HasPrefix(s, "ba:") {
		sp := strings.Split(s, ":")
		for i := range sp {
			fmt.Println(sp[i])
		}
		fmt.Println(fmt.Sprintf("%s:%s", sp[1], sp[2]))
	}
}

func replaceStr() {
	s := "Uncaught SyntaxError: Unexpected token '?'"
	if strings.ContainsAny(s, "'") {
		fmt.Println(strings.Replace(s, "'", "\"", -1))
	}
}
