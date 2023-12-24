package main

import "net/http"

func main() {

	fs := http.FileServer(http.Dir(""))
	http.Handle("/", http.StripPrefix("/", fs))

	http.ListenAndServe(":80", nil)
}
