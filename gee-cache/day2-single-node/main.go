package main

import (
	"fmt"
	"geecache"
	"log"
)

var db = map[string]string{
	"tom": "33",
	"jack": "21",
	"sam": "12",
}

func main()  {
	loadCounts := make(map[string]int, len(db))
	gee := geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			fmt.Println("failed tp get value of tom")
		}
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			fmt.Printf("cache %s miss", k)
		}
	}
	if view, err := gee.Get("unknow"); err == nil {
		fmt.Printf("the value be empty, but %s got", view)
	}
}
