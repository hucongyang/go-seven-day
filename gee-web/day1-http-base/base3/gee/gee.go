package gee

import (
	"fmt"
	"net/http"
)

// 实现的serveHTTP服务

// 向框架用户提供，用来定义路由映射的处理方法
type HandlerFunc func(http.ResponseWriter, *http.Request)

// 定义engine
type Engine struct {
	router map[string]http.HandlerFunc // 路由映射表
}

// 构建 engine对外访问函数
func New() *Engine {
	return &Engine{
		router: make(map[string]http.HandlerFunc),
	}
}

// 实现 net/http server.go 的 Handler 接口
// 解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法，查不到，就返回404
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	routerKey := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[routerKey]; ok {
		handler(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 not found: %s\n", req.URL)
	}
}

// 向 engine router里面增加handler
func (engine *Engine) addRoute(method string, pattern string, handler http.HandlerFunc) {
	routerKey := method + "-" + pattern
	engine.router[routerKey] = handler
}

// 当用户调用 GET 方法时，会将路由和处理方法注册到映射表 router中
func (engine *Engine) GET(pattern string, handler http.HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler http.HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 运行自定义http服务
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
