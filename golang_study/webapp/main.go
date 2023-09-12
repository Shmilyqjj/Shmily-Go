package main

// Middlewares

import (
	"fmt"
)

const maxIndex = 63

type HandlerFunc func(ctx *context)

type context struct {
	HandlersChain []HandlerFunc
	index         int8
}

func (ctx *context) next() {
	if ctx.index < maxIndex {
		ctx.index++
		ctx.HandlersChain[ctx.index](ctx)
	}
}

func (ctx *context) abort() {
	ctx.index = maxIndex
	fmt.Println("abort...")
}

func (ctx *context) use(f HandlerFunc) {
	ctx.HandlersChain = append(ctx.HandlersChain, f)
}

func (ctx *context) get(relativePath string, f HandlerFunc) {
	ctx.HandlersChain = append(ctx.HandlersChain, f)
}

func (ctx *context) run() {
	ctx.HandlersChain[0](ctx)
}

func main() {
	ctx := &context{}
	ctx.use(middleware1)
	ctx.use(middleware2)
	ctx.get("hahahah", logicFunc)
	ctx.run()
}

func middleware1(ctx *context) {
	fmt.Println("middleware1 begin")
	//ctx.abort()
	ctx.next()
	fmt.Println("middleware1 end")
}

func middleware2(ctx *context) {
	fmt.Println("middleware2 begin")
	ctx.next()
	fmt.Println("middleware2 end")
}

func logicFunc(ctx *context) {
	fmt.Println("logicFunc function")
}
