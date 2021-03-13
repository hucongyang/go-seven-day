package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	for index, value := range r.Header {
		fmt.Fprintf(w, "index: %s, value: %s\n", index, value)
	}
}
