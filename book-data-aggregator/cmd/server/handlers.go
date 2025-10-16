package main

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"strconv"
	"sync"
)

type jsonData map[string]interface{}

type result struct {
	TotalBooks int
	EarliestPublicationYear int
	LatestPublicationYear int
	Authors map[string]bool
	Languages map[string]bool
	mu sync.Mutex
}

func calculateStats(res *result, data *jsonData) string {
	if numFound, ok := (*data)["numFound"].(float64); ok{
		res.TotalBooks = int(numFound)
	} else {
		return "ERROR"
	}
	

	books, ok := (*data)["docs"].([]interface{})
	if !ok {
		return "ERROR"
	}

	for _, book := range books {
		if bookData, ok := book.(map[string]interface{}); ok {

			if authorNames, exists := bookData["author_name"]; exists {
				if authorNamesSlice, ok := authorNames.([]interface{}); ok {
					for _, author := range authorNamesSlice {
						if author, ok := author.(string); ok {
							res.Authors[author] = true
						} else {
							continue;
						}
					}
				}
			}

			if languages, exists := bookData["language"]; exists {
				if laguagesSlice, ok := languages.([]interface{}); ok {
					for _, language := range laguagesSlice {
						if language, ok := language.(string); ok {
							res.Languages[language] = true
						} else {
							continue;
						}
					}
				}
			}

			if first_publish_year, exists := bookData["first_publish_year"]; exists {
				if first_publish_year, ok := first_publish_year.(float64); ok {
					firstPublishYearInt := int(first_publish_year)

					if firstPublishYearInt < res.EarliestPublicationYear {
						res.EarliestPublicationYear = firstPublishYearInt
					}
					if firstPublishYearInt > res.LatestPublicationYear {
						res.LatestPublicationYear = firstPublishYearInt
					}
				}
			}
		}
	}
	fmt.Printf("%v\n", *res)

	return "SUCCESS"
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

	res.mu.Lock()
	defer res.mu.Unlock()
	calculateStats(res, &j)
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

	res := &result{
		EarliestPublicationYear: 10000,
		LatestPublicationYear: -1,
		Authors: make(map[string]bool),
		Languages: make(map[string]bool),
	}

	res.dataRequest(title, page)	
	
}