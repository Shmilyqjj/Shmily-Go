package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	md5ValueByteArray := md5sum("hello")
	md5ValueString := fmt.Sprintf("%x", md5ValueByteArray)
	fmt.Println(md5ValueString)
}

func md5sum(input string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hasher.Sum(nil)
}
