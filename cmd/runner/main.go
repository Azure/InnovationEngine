package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello, world!")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
}
