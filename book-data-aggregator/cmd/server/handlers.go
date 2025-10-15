package main

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"strconv"
)

// We need an empty interface type because the json response
// can have multiple types of values
type jsonData interface{}

type result struct {
	TotalBooks int
	EarliestPublicationYear int
	LatestPublicationYear int
	Languages []string
	Authors []string
}

func (res *result) dataRequest(title string, pageNumber int) {
	response, err := http.Get(fmt.Sprintf("https://openlibrary.org/search.json?title=%s&limit=100&page=%d", title, pageNumber))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()

	byteData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var j jsonData

	err = json.Unmarshal(byteData, &j)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v", j)
}


func getData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed.", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()

	title := queryParams.Get("title")
	page, err := strconv.Atoi(queryParams.Get("page"))

	if err != nil {
		fmt.Println(err)
		return
	}

	var res *result

	res.dataRequest(title, page)	
	
}