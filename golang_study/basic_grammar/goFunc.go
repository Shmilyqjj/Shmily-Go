package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	go func() {
		fmt.Println()
	}()

	sli := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	wg := sync.WaitGroup{}
	for k, v := range sli {
		wg.Add(1)

		go func(k, v interface{}) {
			time.Sleep(time.Second)
			fmt.Println(k, v)
			wg.Done()
		}(k, v)
	}
	wg.Wait()
}
