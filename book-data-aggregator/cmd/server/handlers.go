package main

import (
	"fmt"
	"net/http"
)


func getData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed.", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()

	title := queryParams.Get("title")

	fmt.Println(title)
}