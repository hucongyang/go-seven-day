package singleflight

import "sync"

// call 代表正在进行中，或已经结束的请求，使用sync.WaitGroup锁避免重入
type call struct {
	wg 		sync.WaitGroup
	val 	interface{}
	err 	error
}
// 管理不同key的请求(call)
type Group struct {
	mu		sync.Mutex
	m 		map[string]*call
}
// 参数key，参数函数fn
// 针对相同的key，无论Do被调用多少次，函数fn都只会被调用一次，等待fn调用结束了，返回返回值或错误
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// 存在即返回group中数值
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	// 不存在新建并执行fn，把结束存入group
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}