package main

import (
	"fmt"
	"net/http"
	"net/url"
	"io"
	"encoding/json"
	"sync"
	"strconv"
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

type responseDataStruct struct {
	TotalBooks int					`json:"totalBooks"`
	EarliestPublicationYear int		`json:"earliestPublicationYear"`
	LatestPublicationYear int		`json:"latestPublicationYear"`
	Authors []string				`json:"authors"`
	Languages []string				`json:"languages"`
}



func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "./ui/html/home.html")
}



func calculateStats(res *result, data *jsonData) {
	if numFound, ok := (*data)["numFound"].(float64); ok{
		res.TotalBooks = int(numFound)
	} else {
		return 
	}
	

	books, ok := (*data)["docs"].([]interface{})
	if !ok {
		return
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

	return
}

func (res *result) dataRequest(wg *sync.WaitGroup, title string, pageNumber int) {
	defer wg.Done()

	params := url.Values{}

	params.Add("title", title)
	params.Add("page", strconv.Itoa(pageNumber))

	queryString := params.Encode()

	response, err := http.Get(fmt.Sprintf("https://openlibrary.org/search.json?%s", queryString))
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

	res := &result{
		TotalBooks: 0,
		EarliestPublicationYear: 10000,
		LatestPublicationYear: -1,
		Authors: make(map[string]bool),
		Languages: make(map[string]bool),
	}

	var wg sync.WaitGroup
	for i:=1; i<=5; i++ {
		wg.Add(1)
		go res.dataRequest(&wg, title, i)
	}

	wg.Wait()

	var authorList, languageList []string
	
	for author := range res.Authors {
		authorList = append(authorList, author)
	}

	for language := range res.Languages {
		languageList = append(languageList, language)
	}

	responseData :=  responseDataStruct{
		TotalBooks: res.TotalBooks,
		EarliestPublicationYear: res.EarliestPublicationYear,
		LatestPublicationYear: res.LatestPublicationYear,
		Authors: authorList,
		Languages: languageList,
	}

	responseJsonDataBytes, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	
	w.Write(responseJsonDataBytes)
}