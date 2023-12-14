package main

import (
	"math/rand"
	"strings"
	"time"
)

func main() {
	// aa bb cc 中随机选一个
	s := "aa,bb,cc"
	splits := strings.Split(s, ",")
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(splits))
	println(idx)
	println(splits[idx])
}
