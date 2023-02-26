package main

import (
	"container/list"
	"fmt"
	"time"
)

func main() {
	l := list.New()
	for i := 1; i <= 100000; i++ {
		l.PushBack(i)
	}
	batchSize := 1000
	forBatch(l, batchSize)
}

// 按固定批次大小拆分数据并协程处理
func forBatch(l *list.List, batchSize int) {
	length := l.Len()
	batchNum := length / batchSize
	concurrency := 21
	if length%batchSize != 0 {
		batchNum++
	}
	fmt.Printf("Total batch num: %d \n", batchNum)

	con := make(chan int, concurrency)
	var idx int
	for i := l.Front(); i != nil; i = i.Next() {
		idx++
		if idx%batchSize == 0 || idx == length {
			fmt.Printf("达到批次大小执行协程计算 idx: %d \n", idx)
			con <- 1
			go func() {
				defer func() {
					select {
					case <-con:
					default:
					}
				}()
				time.Sleep(time.Duration(3) * time.Second)
			}()
		} else {
			//fmt.Println("continue")
		}
	}

	for {
		select {}
	}

}
