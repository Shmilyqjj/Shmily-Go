package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"unicode/utf8"
)

func main() {
	strConcat()
	splitStr()
	replaceStr()
	invalidUtf8String()
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

/*
*
无效UTF8字符串处理
*/
func invalidUtf8String() {
	s := "\xcf\xcf067" // invalid UTF-8 string
	if !utf8.ValidString(s) {
		logrus.Warn("String is not valid UTF-8, do sanitizeString")
		s = sanitizeString(s)
	} else {
		logrus.Infoln("String is valid UTF-8")
	}
	fmt.Println("Cleaned string:", s)
}

/*
*
处理非UTF8字符串 输出处理替换后的
*/
func sanitizeString(s string) string {
	validStr := make([]rune, 0, len(s))
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError && size == 1 {
			validStr = append(validStr, '?')
		} else {
			validStr = append(validStr, r)
		}
		s = s[size:]
	}
	return string(validStr)
}
