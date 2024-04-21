package main

import "net/http"

// TODO: Add header: Content-Security-Policy: default-src https://developers.redhat.com
func main() {

	fs := http.FileServer(http.Dir(""))
	http.Handle("/", http.StripPrefix("/", fs))

	http.ListenAndServe(":80", nil)
}
