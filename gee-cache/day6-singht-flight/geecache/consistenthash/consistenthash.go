package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 实现一致性哈希算法

// 函数类型
type Hash func(data []byte) uint32

type Map struct {
	hash		Hash		// Hash函数
	replicas	int			// 虚拟节点倍数
	keys 		[]int		// 哈希环 keys
	hashMap		map[int]string	// 虚拟节点和真实节点的映射表 hashMap，键是虚拟节点的哈希值，值是真实节点的名称
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash: 		fn,
		replicas: 	replicas,
		hashMap: 	make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 添加真实节点/机器的 Add() 方法
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 虚拟节点的名称：strconv.Itoa(i) + key
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))	// 计算虚拟节点的哈希值
			m.keys = append(m.keys, hash)	// 添加到哈希环上
			m.hashMap[hash] = key			// 增加虚拟节点和真实节点的映射关系
		}
	}
	sort.Ints(m.keys)		// 哈希环上的哈希值排序
}
// 选择节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// 顺时针找到第一个匹配的虚拟节点的下标idx
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 通过 hashMap 映射得到真实的节点
	return m.hashMap[m.keys[idx%len(m.keys)]]
}