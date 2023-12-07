package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"time"
)

// curl --location 'http://localhost:8080/index' --header 'Content-Type: application/json' --data '{"id": "1", "name": "qjj", "age": 25}'
// curl --location 'http://localhost:8080/index/qjj?id=1&name=qjj&age=25&sex=m&valid=true'

type Service struct {
	ServerConfig *ServerConfig
}

type ServerConfig struct {
	Port                  int `yaml:"port"`
	ReadTimeout           int `yaml:"read_timeout"`
	WriteTimeout          int `yaml:"write_timeout"`
	IdleTimeout           int `yaml:"idle_timeout"`
	MaxConnsPerIP         int `yaml:"max_conns_per_ip"`
	MaxIdleWorkerDuration int `yaml:"max_idle_worker_duration"`
	MaxRequestBodySize    int `yaml:"max_request_body_size"`
}

// Model
type Index struct {
	Id   string `json:"id"   form:"id"`
	Name string `json:"name" form:"name"`
	Age  int    `json:"age"  form:"age"`
}

func (s *Service) Init() error {
	router := fasthttprouter.New()
	// 注册路由
	router.GET("/index", s.Index)
	router.POST("/index", s.Index)
	router.GET("/index/:name", s.Index)

	//启动服务端口
	server := fasthttp.Server{
		Name:                          "rcv_collector_server",
		TLSConfig:                     &tls.Config{InsecureSkipVerify: true},
		ReadTimeout:                   time.Duration(s.ServerConfig.ReadTimeout) * time.Second,
		WriteTimeout:                  time.Duration(s.ServerConfig.WriteTimeout) * time.Second,
		IdleTimeout:                   time.Duration(s.ServerConfig.IdleTimeout) * time.Second,
		MaxConnsPerIP:                 s.ServerConfig.MaxConnsPerIP,
		MaxIdleWorkerDuration:         time.Duration(s.ServerConfig.MaxIdleWorkerDuration) * time.Second,
		MaxRequestBodySize:            s.ServerConfig.MaxRequestBodySize,
		ReadBufferSize:                65535, // ReadBufferSize 限制header大小 避免报错'Too big request header'
		DisableHeaderNamesNormalizing: true,
		DisableKeepalive:              true,
		Handler:                       router.Handler,
	}
	err := server.ListenAndServe(fmt.Sprintf(":%d", s.ServerConfig.Port))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Index(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name")
	ret := "Welcome!"
	if name != nil {
		ret = fmt.Sprintf("Welcome %s!", name.(string))
	}

	var r Index
	if ctx.IsPost() {
		// Post请求 获取参数
		err := json.Unmarshal(ctx.Request.Body(), &r)
		if err != nil {
			ctx.Error(err.Error(), 1)
			return
		}
		ret = fmt.Sprintf("%s %s %s %d", ret, r.Id, r.Name, r.Age)
	} else if ctx.IsGet() {
		// Get请求 获取参数
		args := ctx.QueryArgs()
		// 遍历获取参数
		args.VisitAll(func(k, v []byte) {
			switch string(k) {
			case "id":
				ret = fmt.Sprintf("%s id=%s ", ret, string(v))
			case "name":
				ret = fmt.Sprintf("%s name=%s ", ret, string(v))
			default:
				ret = fmt.Sprintf("%s [%s=%s] ", ret, k, string(v))
			}
		})
		// 直接获取参数
		id := string(args.Peek("id"))
		n := string(args.Peek("name"))
		valid := args.GetBool("valid")
		age, _ := args.GetUint("age")
		ret = fmt.Sprintf("%s %s %s %d %v", ret, id, n, age, valid)
		ret = fmt.Sprintf("Welcome %s %s!", name.(string), ret)
	}
	_, _ = ctx.WriteString(ret)
}

func main() {
	s := &Service{ServerConfig: &ServerConfig{
		Port:               8080,
		MaxRequestBodySize: 10485760,
	},
	}
	err := s.Init()
	if err != nil {
		panic(err)
	}
}
