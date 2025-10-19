# Book Data Aggregator

A small Golang project created to learn about handling goroutines and mutex in Golang.

This project uses OpenLibrary api to fetch and aggregate book data.


It takes book title as an input an outputs some basic data.

It outputs the data of the first 500 books found. The data it outputs are,

- Total number of books found in with similar title.

- The earliest publication year among the books

- The lateset publication year among the books

- List of all authors of the books

- List of languages of the books


To run the project use your terminal to navigate to the directiory with the "go.mod" file.
Then run the following command:

- go run ./cmd/server

The command will start the server at 127.0.0.1:8000
