package main

import (
	"fmt"
	"strings"
)

func main() {
	strConcat()
	splitStr()
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
	}

}
