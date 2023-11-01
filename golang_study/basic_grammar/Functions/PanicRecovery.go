package main

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

// 避免panic, 使程序恢复正常执行

// recover()
// 1、内建函数
// 2、用来控制一个goroutine的panicking行为，捕获panic，从而影响应用的行为
// 3、一般的调用建议
//   a). 在defer函数中，通过recover来终止一个goroutine的panicking过程，从而恢复正常代码的执行
//   b). 可以获取通过panic传递的error

func main() {
	WithRecoverAndHandle(testPanicFunc, func(i interface{}) {
		logrus.Error(i)
	})

	panicAutoRecoveryFunc()

	println("finished.")

}

func testPanicFunc() {
	panic("哈哈,你小子恐慌了.")
}

func panicAutoRecoveryFunc() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorln(err)
		}
	}()
	time.Sleep(2 * time.Second)
	panic("哈哈,你小子又恐慌了.")
}

type PanicHandle func(interface{})

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	logrus.Errorf("==> %s\n", string(buf[:n]))
}

func Recover() {
	if err := recover(); err != nil {
		// 这里的err其实就是panic传入的内容，55
		logrus.Errorln(err)
		PrintStack()
	}
}

func WithRecover(fn func()) {
	defer Recover()
	fn()
}

func WithRecoverAndHandle(fn func(), handle PanicHandle) {
	defer RecoverWithHandle(handle)
	fn()
}

func RecoverWithHandle(handle PanicHandle) {
	if err := recover(); err != nil {
		handle(err)
	}
}
