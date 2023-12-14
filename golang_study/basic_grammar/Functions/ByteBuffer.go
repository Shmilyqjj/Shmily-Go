package main

import "bytes"

func main() {
	dataBuffer := bytes.Buffer{}
	s1b := []byte("aaa")
	s2b := []byte("bbb")
	dataBuffer.Write(s1b)
	dataBuffer.Write([]byte("\n"))
	dataBuffer.Write(s2b)
	println(dataBuffer.String())
	println(string(dataBuffer.Bytes()))
	// 使用content-type application/octet-stream 可以模拟发送文件 将字节流以文件形式发给接口
}
