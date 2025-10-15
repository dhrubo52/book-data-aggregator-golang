package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/get-data", getData)

	fmt.Println("Starting Server at 127.0.0.1:5000")

	err := http.ListenAndServe("127.0.0.1:5000", mux)
	if err!=nil {
		fmt.Println(err)
	}
}