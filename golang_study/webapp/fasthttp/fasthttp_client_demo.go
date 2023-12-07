package main

import (
	"crypto/tls"
	"fmt"
	"github.com/valyala/fasthttp"
)

type ClientConfig struct {
	Port                  int `yaml:"port"`
	ReadTimeout           int `yaml:"read_timeout"`
	WriteTimeout          int `yaml:"write_timeout"`
	IdleTimeout           int `yaml:"idle_timeout"`
	MaxConnsPerIP         int `yaml:"max_conns_per_ip"`
	MaxIdleWorkerDuration int `yaml:"max_idle_worker_duration"`
	MaxRequestBodySize    int `yaml:"max_request_body_size"`
}

func main() {
	get()
	println("============================")
	post()
	println("============================")
	complexRequest()
	println("============================")
	highPerformanceRequest()
	println("============================")
	originalClientRequest()
}

func get() {
	url := `http://httpbin.org/get`

	status, resp, err := fasthttp.Get(nil, url)
	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return
	}

	if status != fasthttp.StatusOK {
		fmt.Println("请求没有成功:", status)
		return
	}

	fmt.Println(string(resp))
}

func post() {
	url := `http://httpbin.org/post?key=123`

	// 填充表单，类似于net/url
	args := &fasthttp.Args{}
	args.Add("name", "test")
	args.Add("age", "18")

	status, resp, err := fasthttp.Post(nil, url, args)
	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return
	}

	if status != fasthttp.StatusOK {
		fmt.Println("请求没有成功:", status)
		return
	}

	fmt.Println(string(resp))
}

func complexRequest() {
	url := `http://httpbin.org/post?key=123`

	req := &fasthttp.Request{}
	req.SetRequestURI(url)

	requestBody := []byte(`{"request":"test"}`)
	req.SetBody(requestBody)

	// 默认是application/x-www-form-urlencoded
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")

	resp := &fasthttp.Response{}

	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		fmt.Println("请求失败:", err.Error())
		return
	}

	b := resp.Body()

	fmt.Println("result:\r\n", string(b))
}

func highPerformanceRequest() {
	url := `http://httpbin.org/post?key=123`

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 用完需要释放资源

	// 默认是application/x-www-form-urlencoded
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")

	req.SetRequestURI(url)

	requestBody := []byte(`{"request":"test"}`)
	req.SetBody(requestBody)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源

	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Println("请求失败:", err.Error())
		return
	}

	b := resp.Body()

	fmt.Println("result:\r\n", string(b))
}

func originalClientRequest() {
	url := `http://httpbin.org/post?key=123`

	client := fasthttp.Client{
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 用完需要释放资源
	// 默认是application/x-www-form-urlencoded
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")

	req.SetRequestURI(url)

	requestBody := []byte(`{"request":"test"}`)
	req.SetBody(requestBody)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源

	if err := client.Do(req, resp); err != nil {
		fmt.Println("请求失败:", err.Error())
		return
	}

	b := resp.Body()

	fmt.Println("result:\r\n", string(b))
}
