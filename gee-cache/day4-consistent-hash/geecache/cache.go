package geecache

import (
	"geecache/lru"
	"sync"
)

// 并发控制
// 封装lru接口，支持并发操作
// 实例化lru，封装get和add方法，并添加互斥锁mu
type cache struct {
	mu			sync.Mutex
	lru 		*lru.Cache
	cacheBytes	int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}