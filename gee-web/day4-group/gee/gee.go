package gee

import (
	"log"
	"net/http"
)

// 相当于 handlerFunc(w http.ResponseWriter, req *http.Request)
type HandlerFunc func(c *Context)

type (
	// 路由分组
	RouterGroup struct {
		prefix 			string			// 前缀
		middlewares 	[]HandlerFunc	// 中间件
		parent			*RouterGroup	// 父级路径
		engine 			*Engine
	}

	// engine作为最顶层的分组，具有RouterGroup所有的能力
	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup	// store all groups
	}
)

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	// 实现了路由的映射, engine从某种意义上继承了 RouterGroup的所有属性和方法
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}