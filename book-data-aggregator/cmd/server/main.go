package main

import (
	"fmt"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/get-data", getData)

	fmt.Printf("Starting Server at 127.0.0.1:8000\n")
	fmt.Printf("Visit: 127.0.0.1:8000\n")

	err := http.ListenAndServe("127.0.0.1:8000", mux)
	if err!=nil {
		fmt.Println(err)
	}
}