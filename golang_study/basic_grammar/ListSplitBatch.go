package main

import (
	"fmt"
	"math"
	"sync"
)

func main() {
	// 创建一个包含 100 个元素的列表
	list := make([]int, 100)
	for i := 0; i < 100; i++ {
		list[i] = i
	}

	// 设置批次大小为 10
	batchSize := 10

	// 计算批次数量
	batchCount := int(math.Ceil(float64(len(list)) / float64(batchSize)))

	// 创建一个 WaitGroup，用于等待所有 goroutine 完成
	var wg sync.WaitGroup
	wg.Add(batchCount)

	// 创建一个 channel，用于传递数据
	ch := make(chan []int)

	// 启动 goroutine 处理数据
	for i := 0; i < batchCount; i++ {
		// 计算批次的起始和结束下标
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(list) {
			end = len(list)
		}

		// 将批次数据发送到 channel 中
		go func(data []int) {
			ch <- data
		}(list[start:end])

		// 启动一个 goroutine 处理每个批次的数据
		go func() {
			defer wg.Done()

			for data := range ch {
				// 处理数据
				for _, value := range data {
					// TODO: 处理数据的逻辑
					fmt.Println(value)
				}
			}
		}()
	}

	// 等待所有 goroutine 完成
	wg.Wait()
}
