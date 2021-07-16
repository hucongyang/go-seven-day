package main

import (
	"gee"
	"net/http"
)

func main()  {
	r := gee.Default()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})
	r.Run(":9999")
}