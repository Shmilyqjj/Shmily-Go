package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	channel()

	forSelectChan()

	waitForSignal()

	forSelectChanWg()

	// 定时触发和消息触发同时进行
	tickerAndMessageIn()
}

func channel() {
	// 信道是带有类型的管道，你可以通过它用信道操作符 <- 来发送或者接收值。
	// 默认情况下，发送和接收操作在另一端准备好之前都会阻塞。这使得goroutine可以在没有显式的锁或竞态变量的情况下进行同步。
	// 适goroutine间的通信
	s := []int{7, 2, 8, -9, 4, 0}

	// 创建chan信道 类型int
	//使用两个goroutine并行进行求和
	ch := make(chan int)
	go chanSum(s[:len(s)/2], ch)
	go chanSum(s[len(s)/2:], ch)
	x, y := <-ch, <-ch // 从 c 中接收
	fmt.Println(x, y, x+y)

	// 创建带缓冲的信道chan
	// 仅当信道的缓冲区填满后，向其发送数据时才会阻塞。当缓冲区为空时，接受方会阻塞。
	ch = make(chan int, 2)
	ch <- 6
	ch <- 8
	println(<-ch)
	v, ok := <-ch
	println(v, ok)
	// 发送者可通过 close 关闭一个信道来表示没有需要发送的值了。接收者可以通过为接收表达式分配第二个参数来测试信道是否被关闭
	close(ch)
	v, ok = <-ch
	println(v, ok)
}

func chanSum(ints []int, ch chan int) {
	sum := 0
	for _, intVal := range ints {
		sum += intVal
	}
	ch <- sum
}

func forSelectChan() {
	// select 语句使一个 Go 程可以等待多个通信操作。
	// select 通常与 for 一起使用，以创建一个无限循环，该循环等待多个通道操作中的一个完成
	// select 通常用于多路复用不同的通道操作，它需要在一个循环中不断尝试这些操作，直到其中一个成功。

	quit := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			println("sleep 1s")
			time.Sleep(1 * time.Second)
		}
		quit <- 0
	}()
	// 等待协程操作完成
	for {
		select {
		case <-quit:
			println("quit")
			return
		default:
			// select所有其他分支未执行时 执行该分支
			// 为了在尝试发送或者接收时不发生阻塞
			println("waiting")
			time.Sleep(3 * time.Second)
		}
	}
}

// =======================================

func forSelectChanWg() {
	stop := make(chan bool)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(stop <-chan bool) {
			defer wg.Done()
			consumer(stop)
		}(stop)
	}
	waitForSignal()
	close(stop)
	fmt.Println("stopping all jobs!")
	wg.Wait()
}

func consumer(stop <-chan bool) {
	for {
		select {
		case <-stop:
			fmt.Println("exit sub goroutine")
			return
		default:
			fmt.Println("running...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func waitForSignal() {
	// 在主程序运行期间等待中断或终止信号的到来。一旦信号被捕获，主程序通常会执行一些清理操作，关闭资源，并终止程序的执行。
	// 这使得程序可以在接收到信号时安全地停止，而不会造成数据损失或不一致性状态。
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt)
	signal.Notify(sigs, syscall.SIGTERM)
	<-sigs
}

// =======================================

func tickerAndMessageIn() {
	// 消息写入buffer 等待定时发送
	ticker := time.NewTicker(5 * time.Second)
	messagesChan := make(chan string, 100)
	dataBuffer := bytes.Buffer{}
	defer close(messagesChan)
	defer ticker.Stop()

	go func() {
		// 产生消息
		i := 0
		for {
			i++
			messagesChan <- fmt.Sprintf("message_%d", i)
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		case message := <-messagesChan:
			dataBuffer.Write([]byte(message))
			dataBuffer.Write([]byte("\n"))
		case <-ticker.C:
			if dataBuffer.Len() != 0 {
				println("send: \n" + dataBuffer.String())
				dataBuffer.Reset()
			}
		}

	}
}
