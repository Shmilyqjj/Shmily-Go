package main

import (
	"fmt"
	"time"
)

// 定时器 可以用于定时执行操作或定时动态更新配置文件
func main() {
	ticker := time.NewTicker(time.Second * 1) //创建一个周期性定时器
	defer ticker.Stop()

	// 协程启动 避免阻塞
	go func() {
		for {
			select {
			case <-ticker.C:
				// reload config
				fmt.Println("Refresh config_0 ...")
			}
		}
	}()

	// 非协程启动 阻塞
	for {
		select {
		case <-ticker.C:
			// reload config
			fmt.Println("Refresh config_1 ...")
		}
	}

}
