package main

// Middlewares

import (
	"fmt"
)

const maxIndex = 63

type HandlerFunc func(ctx *contextM)

type contextM struct {
	HandlersChain []HandlerFunc
	index         int8
}

func (ctx *contextM) next() {
	if ctx.index < maxIndex {
		ctx.index++
		ctx.HandlersChain[ctx.index](ctx)
	}
}

func (ctx *contextM) abort() {
	ctx.index = maxIndex
	fmt.Println("abort...")
}

func (ctx *contextM) use(f HandlerFunc) {
	ctx.HandlersChain = append(ctx.HandlersChain, f)
}

func (ctx *contextM) get(relativePath string, f HandlerFunc) {
	ctx.HandlersChain = append(ctx.HandlersChain, f)
}

func (ctx *contextM) run() {
	ctx.HandlersChain[0](ctx)
}

func main() {
	ctx := &contextM{}
	ctx.use(middleware1)
	ctx.use(middleware2)
	ctx.get("hahahah", logicFunc)
	ctx.run()
}

func middleware1(ctx *contextM) {
	fmt.Println("middleware1 begin")
	//ctx.abort()
	ctx.next()
	fmt.Println("middleware1 end")
}

func middleware2(ctx *contextM) {
	fmt.Println("middleware2 begin")
	ctx.next()
	fmt.Println("middleware2 end")
}

func logicFunc(ctx *contextM) {
	fmt.Println("logicFunc function")
}
