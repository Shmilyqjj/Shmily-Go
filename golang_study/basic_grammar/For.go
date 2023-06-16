package main

import (
	"container/list"
	"fmt"
	"strings"
	"time"
)

func main() {
	l := list.New()
	for i := 1; i <= 100000; i++ {
		l.PushBack(i)
	}
	//batchSize := 1000
	//forBatch(l, batchSize)

	forRangeListMap()
	fmt.Println(map2String(map[string]string{"name": "John", "age": "30"}))
	fmt.Println(string2Map("k1=v1,k2=v2,k3=v3"))
	fmt.Println(splitAndDistinct("1=1,2=2,3=3,3=3,3=3"))
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

func forRangeListMap() {
	Data := []map[string]string{
		{"name": "John", "age": "30"},
		{"name": "Alice", "age": "25"},
		{"name": "Bob", "age": "35"},
	}

	// 遍历切片中的每个 map
	for _, m := range Data {
		// 遍历 map 中的键值对
		for key, value := range m {
			fmt.Printf("键：%s，值：%s\n", key, value)
		}
		fmt.Println("--------------------")
	}
}

func map2String(m map[string]string) string {
	l := make([]string, len(m))
	i := 0
	for k, v := range m {
		l[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}
	return strings.Join(l, ",")
}

func string2Map(s string) map[string]string {
	// k1=v1,k2=v2转map
	sp := strings.Split(s, ",")
	m := make(map[string]string, len(sp))
	for _, pair := range sp {
		keyValue := strings.Split(pair, "=")
		if len(keyValue) == 2 {
			key := keyValue[0]
			value := keyValue[1]
			m[key] = value
		}
	}
	return m
}

func splitAndDistinct(s string) []string {
	// split结果去重
	sp := strings.Split(s, ",")
	set := make(map[string]struct{})
	var uniqueElements []string
	for _, element := range sp {
		if _, exists := set[element]; !exists {
			set[element] = struct{}{}
			uniqueElements = append(uniqueElements, element)
		}
	}
	return uniqueElements
}
