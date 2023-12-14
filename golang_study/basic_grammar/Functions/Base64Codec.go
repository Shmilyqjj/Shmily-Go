package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	userName := "admin"
	password := "123456"
	up := []byte(fmt.Sprintf("%s:%s", userName, password))
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(up)))
	base64.StdEncoding.Encode(encoded, up)
	println(fmt.Sprintf("Basic %s", encoded))
}
