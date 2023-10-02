package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/basic_auth"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// References: https://github.com/cloudwego/hertz-examples/tree/main

func main() {
	h := server.Default(
		server.WithHostPorts("127.0.0.1:8080"),
		server.WithMaxRequestBodySize(20<<20),
		server.WithTransport(standard.NewTransporter),
	)

	h.GET("/hello", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "Hello hertz!")
	})

	RegisterRoute(h)
	RegisterGroupRoute(h)
	RegisterGroupRouteWithMiddleware(h)

	h.Spin()
}

func RegisterRoute(h *server.Hertz) {
	h.StaticFS("/", &app.FS{Root: "./", GenerateIndexPages: true})

	h.GET("/get", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "get")
	})
	h.POST("/post", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "post")
	})
	h.PUT("/put", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "put")
	})
	h.DELETE("/delete", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "delete")
	})
	h.PATCH("/patch", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "patch")
	})
	h.HEAD("/head", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "head")
	})
	h.OPTIONS("/options", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "options")
	})
	h.Any("/ping_any", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "any")
	})
	h.Handle("LOAD", "/load", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "load")
	})
}

func RegisterGroupRoute(h *server.Hertz) {
	// Simple group: v1
	v1 := h.Group("/v1")
	{
		// loginEndpoint is a handler func
		v1.GET("/get", func(ctx context.Context, c *app.RequestContext) {
			c.String(consts.StatusOK, "get")
		})
		v1.POST("/post", func(ctx context.Context, c *app.RequestContext) {
			c.String(consts.StatusOK, "post")
		})
	}

	// Simple group: v2
	v2 := h.Group("/v2")
	{
		v2.PUT("/put", func(ctx context.Context, c *app.RequestContext) {
			c.String(consts.StatusOK, "put")
		})
		v2.DELETE("/delete", func(ctx context.Context, c *app.RequestContext) {
			c.String(consts.StatusOK, "delete")
		})
	}
}

func RegisterGroupRouteWithMiddleware(h *server.Hertz) {
	example1 := h.Group("/example1", basic_auth.BasicAuth(map[string]string{"test": "test"}))
	example1.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(consts.StatusOK, "pong")
	})

	example2 := h.Group("/example2")
	example2.Use(basic_auth.BasicAuth(map[string]string{"test": "test"}))
	example2.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(consts.StatusOK, "pong")
	})
}
