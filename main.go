package main

import (
	"fmt"
	"log"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}

func main() {

	http.HandleFunc("/", root)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
