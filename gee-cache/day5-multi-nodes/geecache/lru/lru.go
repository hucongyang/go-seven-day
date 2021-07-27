package lru

import (
	"container/list"
)

/**
缓存淘汰策略：最近最少使用
lru：如果数据最近被访问过，那么将来被访问的概率会更高
原理：维护一个队列，如果某条记录被访问了，则移动到队尾，那么队首则是最近最少访问的数据，淘汰该条记录即可
 */

type Cache struct {
	maxBytes 	int64		// 允许使用的最大内存
	nbytes 		int64		// 当前已使用的内存
	ll 			*list.List
	cache		map[string]*list.Element
	OnEvicted	func(key string, value Value)	// 某条记录被移除时的回调函数
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 双向链表节点的数据类型
type entry struct {
	key 	string
	value	Value
}

type Value interface {
	Len() int		// 用于返回值所占用的内存大小
}

// 查找功能，第一步从字典中找到对应的双向链表的节点，第二步，将该节点移动到队尾
// 如果键对应的链表节点存在，则将对应节点移动到队尾，并返回查找到的值
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 删除功能: 缓存淘汰，即移除最近最少访问的节点（队首）
// c.ll.Back() 取到队首节点，从链表中删除
// delete(c.cache, kv,key) 从字典中 c.cache 删除该节点的映射关系
// 如果回调函数 OnEvicted 不为nil，则调用回调函数
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 新增/修改
// 如果键存在，则更新对应节点的值，并将该节点移到队尾
// 不存在则是新增场景，首先队尾添加新节点，并在字典中添加key和节点的映射关系
// 更新 c.nbytes 如果超过了设定的最大值 c.maxBytes 则移除最少访问的节点
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(len(kv.key)) + int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// 用来获取添加了多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}