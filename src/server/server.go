package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/search", Search)
	http.HandleFunc("/results", Results)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
