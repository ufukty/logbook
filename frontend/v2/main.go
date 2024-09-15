package main

import (
	"net/http"
	"os"
	"path/filepath"
)

// TODO: Add header: Content-Security-Policy: default-src https://developers.redhat.com
func main() {
	root := filepath.Join(os.Getenv("WORKSPACE"), "frontend/cmd/spa/v2")
	fs := http.FileServer(http.Dir(root))
	http.Handle("/", http.StripPrefix("/", fs))
	http.ListenAndServe(":80", nil)
}
