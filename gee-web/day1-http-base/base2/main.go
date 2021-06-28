package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct {}

// 第一个参数可以构造针对该请求的响应，第二个参数包含了该HTTP请求的所有的信息
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.PATH = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 not found: %s\n", req.URL)
	}
}

func main()  {
	engine := new(Engine)
	// 第一个参数是地址，第二个参数是代表处理所有的HTTP请求的实例,nil代表使用标准库中的实例处理
	log.Fatal(http.ListenAndServe(":9999", engine))
}